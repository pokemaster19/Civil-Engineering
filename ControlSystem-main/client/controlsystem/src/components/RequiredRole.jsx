import { useAuth } from "../context/AuthContext";
import {ROLES} from "../constants/Roles"

export const RequireRole = ({ allowedRoles, children }) => {
    const { roleId, isLoading } = useAuth();

    if (isLoading) return null;

    if (roleId == ROLES.ADMIN || roleId == ROLES.SUPERADMIN){
        return children;
    }

    if (!allowedRoles.includes(roleId)) {
        return null;
    }

    return children;
};