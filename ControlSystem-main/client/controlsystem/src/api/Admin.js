import api from './axiosInstance'

export const registerNewUser = async (firstName, middleName, lastName, email, role) =>{
    try{
        const response = await api.axiosInstance.post('/admin/register', {
            first_name: firstName,
            middle_name: middleName,
            last_name: lastName,
            email: email,
            role: parseInt(role, 10),
        });
        return response.data;

    }
    catch(err){

    }
}

export const editUser = async (userId, firstName, middleName, lastName, role, isEnabled) => {
    try{
        const response = await api.axiosInstance.patch(`/admin/edit-user/${userId}`, {
            first_name: firstName,
            middle_name: middleName,
            last_name: lastName,
            role: role,
            is_enabled: isEnabled,
        });
        return response.data;
    }
    catch (err){

    }
}


export const getAllUsers = async ({ page = 1, email = '', role = '', isEnabled = '' } = {}) => {
    try {
        const params = new URLSearchParams();

        params.append('page', page);
        if (email) params.append('email', email);
        if (role) params.append('role', role);
        if (isEnabled !== '') params.append('is_enabled', isEnabled);

        const response = await api.axiosInstance.get(`/admin/get-users?${params.toString()}`);

        return {
            users: Array.isArray(response.data.users) ? response.data.users : [],
            pagination: response.data.pagination || { limit: 10, page: 1, total: 0, totalPages: 1 }
        };
    } catch (err) {
        console.error('Ошибка в getAllUsers:', err);
        return {
            users: [],
            pagination: { limit: 10, page: 1, total: 0, totalPages: 1 }
        };
    }
};
