import { createContext, useState, useContext, useEffect, useCallback } from 'react';
import api from '../api/axiosInstance';

const AuthContext = createContext(null);

export const AuthProvider = ({ children }) => {
    const [accessToken, setAccessToken] = useState(null);
    const [userId, setUserId] = useState(null);
    const [roleId, setRoleId] = useState(null);
    const [isLoading, setIsLoading] = useState(true);

    const setAuthState = useCallback((newAccessToken, newUserId, newRoleId) => {
        setAccessToken(newAccessToken);
        setUserId(newUserId);
        setRoleId(newRoleId);
        api.setAxiosAuthToken(newAccessToken);
    }, []);

    const clearAuthState = useCallback(() => {
        setAccessToken(null);
        setUserId(null);
        setRoleId(null);
        api.clearAxiosAuthToken();
    }, []);

    const login = useCallback((newAccessToken, newUserId, newRoleId) => {
        console.log("AuthContext: Logging in...");
        setAuthState(newAccessToken, newUserId, newRoleId);
    }, [setAuthState]);

    const logout = useCallback(async (triggeredByRefreshFailure = false) => {
        console.log("AuthContext: Logging out", triggeredByRefreshFailure ? "(due to refresh failure)" : "");
        
        if (!triggeredByRefreshFailure) {
            try {
                await api.authAxiosInstance.post('auth/logout');
                console.log("AuthContext: Logout request sent to backend");
            } catch (error) {
                console.error("AuthContext: Error during backend logout:", error);
            }
        }
        clearAuthState();
    }, [clearAuthState]);

    const refreshAccessToken = useCallback(async () => {
        console.log("AuthContext: Attempting to refresh token...");
        try {
            const response = await api.authAxiosInstance.post('auth/refresh');
            const { token: newAccessToken, user_id: newUserId, role: newRoleId } = response.data;
            console.log("AuthContext: Token refreshed successfully.");
            setAuthState(newAccessToken, newUserId || userId, newRoleId || roleId);
            return newAccessToken;
        } catch (error) {
            console.error("AuthContext: Failed to refresh token:", error);
            clearAuthState();
            throw error;
        }
    }, [setAuthState, clearAuthState, userId, roleId]);

    useEffect(() => {
        const handleTokenRefreshed = (event) => {
            const { newAccessToken } = event.detail;
            console.log("AuthContext: Received token-refreshed event. Updating state.");
            setAccessToken(newAccessToken);
        };

        window.addEventListener('token-refreshed', handleTokenRefreshed);

        return () => {
            window.removeEventListener('token-refreshed', handleTokenRefreshed);
        };
    }, []);

    useEffect(() => {
        const handleRefreshFailure = () => {
            console.log("AuthContext: Received auth-refresh-failed event. Logging out.");
            logout(true);
        };
        window.addEventListener('auth-refresh-failed', handleRefreshFailure);
        return () => {
            window.removeEventListener('auth-refresh-failed', handleRefreshFailure);
        };
    }, [logout]);
    
    useEffect(() => {
        const attemptRefreshOnLoad = async () => {
            console.log("AuthContext: Attempting initial refresh on load...");
            try {
                await refreshAccessToken();
                console.log("AuthContext: Initial refresh successful.");
            } catch (error) {
                console.log("AuthContext: No valid refresh token found on load or refresh failed.");
            } finally {
                setIsLoading(false);
                console.log("AuthContext: Initial loading finished.");
            }
        };
        attemptRefreshOnLoad();
    }, [refreshAccessToken]);

    const value = {
        accessToken,
        userId,
        roleId,
        isAuthenticated: !!accessToken,
        isLoading,
        login,
        logout,
        refreshAccessToken,
    };

    return (
        <AuthContext.Provider value={value}>
            {children}
        </AuthContext.Provider>
    );
};

export const useAuth = () => {
    const context = useContext(AuthContext);
    if (!context) {
        throw new Error('useAuth must be used within an AuthProvider');
    }
    return context;
};