import { Navigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import { ROLES } from "../constants/Roles";

export const AdminRoute = ({ allowedRoles, children }) => {
    const { roleId, isAuthenticated, isLoading } = useAuth();

    if (isLoading) return <div>Загрузка...</div>;

    if (!isAuthenticated) {
        return <Navigate to="/login" replace />;
    }

    if (roleId >= ROLES.ADMIN) {
        return children;
    }

    else {
        return <Navigate to="/" replace />;
    }

};