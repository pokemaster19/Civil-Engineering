import axios from 'axios';

const API_BASE_URL = process.env.REACT_APP_BACKEND_URL;

const axiosInstance = axios.create({
    baseURL: API_BASE_URL,
    withCredentials: true,
});

const authAxiosInstance = axios.create({
    baseURL: API_BASE_URL,
    withCredentials: true,
});

let isRefreshing = false;
let failedQueue = [];

const processQueue = (error, token = null) => {
    failedQueue.forEach(prom => {
        if (error) {
            prom.reject(error);
        } else {
            prom.resolve(token);
        }
    });
    failedQueue = [];
};

axiosInstance.interceptors.response.use(
    (response) => {
        return response;
    },
    async (error) => {
        const originalRequest = error.config;

        if (error.response?.status === 401 && !originalRequest._retry) {
            if (isRefreshing) {
                return new Promise(function(resolve, reject) {
                    failedQueue.push({ resolve, reject });
                }).then(token => {
                    originalRequest.headers['Authorization'] = 'Bearer ' + token;
                    return axiosInstance(originalRequest);
                }).catch(err => {
                    return Promise.reject(err);
                });
            }

            originalRequest._retry = true;
            isRefreshing = true;

            try {
                console.log("Interceptor: Detected 401, attempting refresh...");
                const refreshResponse = await authAxiosInstance.post('auth/refresh');
                const { token: newAccessToken } = refreshResponse.data;
                console.log("Interceptor: Refresh successful, got new token.");

                window.dispatchEvent(new CustomEvent('token-refreshed', { detail: { newAccessToken } }));

                axiosInstance.defaults.headers.common['Authorization'] = 'Bearer ' + newAccessToken;
                originalRequest.headers['Authorization'] = 'Bearer ' + newAccessToken;

                processQueue(null, newAccessToken);

                return axiosInstance(originalRequest);

            } catch (refreshError) {
                console.error("Interceptor: Refresh token failed:", refreshError);
                processQueue(refreshError, null);
                window.dispatchEvent(new CustomEvent('auth-refresh-failed'));
                return Promise.reject(refreshError);
            } finally {
                isRefreshing = false;
            }
        }

        return Promise.reject(error);
    }
);

const setAxiosAuthToken = (token) => {
    axiosInstance.defaults.headers.common['Authorization'] = token ? `Bearer ${token}` : undefined;
};

const clearAxiosAuthToken = () => {
    delete axiosInstance.defaults.headers.common['Authorization'];
};

const api = {
    axiosInstance,
    authAxiosInstance,
    setAxiosAuthToken,
    clearAxiosAuthToken,
};

export default api;