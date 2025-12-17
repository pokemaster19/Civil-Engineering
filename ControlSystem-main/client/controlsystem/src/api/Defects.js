import api from './axiosInstance'

export const fetchAllDefects = async (projectId, { page = 1} = {}, search = '') => {
    try{
        const response = await api.axiosInstance.get(`/defects/?projectId=${projectId}&page=${page}&search=${search}`);
        return response.data;
    }
    catch (err){
        console.error("API Error fetchAllDefects:", err.response || err.message || err);
        throw new Error(err.response?.data?.error || err.message || "Failed to fetch projects information");
    }
}

export const fetchDefectById = async (defectId) =>{
    try{
        const response = await api.axiosInstance.get(`/defects/${defectId}`);
        return response.data;
    }
    catch (err){
        console.error(err.response || err.message || err)
        throw new Error(err.response?.data?.error || err.message)
    }
}

export const createDefect = async (project_id, title, description) => {
    try{
        const response = await api.axiosInstance.post(`/defects/`, {
            title: title,
            description: description,
            project_id: parseInt(project_id, 10),
        })
        return response;
    }
    catch (err){
        console.error(err.response || err.message || err)
        throw new Error(err.response?.data?.error || err.message)
    }
}

export const editDefect = async (defectId, title, description, priority, status) => {
    try{
        const response = await api.axiosInstance.patch(`/defects/${defectId}`, {
            title: title,
            description: description,
            priority: parseInt(priority, 10),
            status: parseInt(status, 10),
        });
        return response.data;
    }
    catch (err){
        console.error(err.response || err.message || err)
        throw new Error(err.response?.data?.error || err.message)
    }
}

export const leaveComment = async (defectId, comment) => {
    try{
        const response = await api.axiosInstance.post(`/defects/${defectId}/comments`, {
            content: comment,
        });
        return response.data;
    }
    catch (err){
        console.error(err.response || err.message || err)
        throw new Error(err.response?.data?.error || err.message)
    }
}

export const fetchComments = async (defectId, page = 1, limit = 10) => {
    try{
        const response = await api.axiosInstance.get(`/defects/${defectId}/comments?page=${page}&limit=${limit}`);
        return response.data;
    }
    catch (err){
        
    }
}


