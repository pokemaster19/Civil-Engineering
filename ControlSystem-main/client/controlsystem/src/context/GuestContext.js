import { Navigate } from "react-router-dom";
import { useAuth } from "./AuthContext";

export const GuestRoute = ({ children }) => {
    const { isAuthenticated, isLoading } = useAuth();

    if (isLoading) {
        return <div>Загрузка...</div>; 
    }
    if (isAuthenticated) {
        return <Navigate to="/" replace />;
    }

    return children;
};