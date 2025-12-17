import api from "./axiosInstance";

export const logout = async () => {
    try{
        const response = await api.axiosInstance.post(`/auth/logout`);
        return response;
    }
    catch (err){

    }
}