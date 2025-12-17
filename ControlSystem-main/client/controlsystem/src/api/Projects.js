import api from './axiosInstance'

export const fetchAllProjects = async (page = 1, search = '') => {
    try {
        const limit = 10;
        const response = await api.axiosInstance.get(`/projects/?page=${page}&limit=${limit}&search=${search}`);
        return response.data;
    } catch (err) {
        console.error("API Error fetchAllProjects:", err.response || err.message || err);
        throw new Error(err.response?.data?.error || err.message || "Failed to fetch projects information");
    }
};

export const fetchProjectById = async (projectId) => {
    try {
        const response = await api.axiosInstance.get(`/projects/${projectId}`);
        return response.data;
    }
    catch (err){
        console.error(err.response || err.message || err)
        throw new Error(err.response?.data?.error || err.message)
    }
}

export const createProject = async (title, description) => {
    try{
        const response = await api.axiosInstance.post('/projects/', {
            name: title,
            description: description,
        })
        return response.data;
    }
    catch (err){
        throw new Error(`Ошибка при создании проекта: ${err.message}`);
    }
}

export const editProject = async (projectId, projectName, description, status) => {
    try{
        const response = await api.axiosInstance.patch(`/projects/${projectId}`, {
            name: projectName,
            description: description,
            status: parseInt(status, 10),
        });
        return response.data;
    }
    catch (err){
        console.error(err.response || err.message || err)
        throw new Error(err.response?.data?.error || err.message)
    }
}

export const assignEngineer = async (projectId, engineerId) => {
    try{
        const response = await api.axiosInstance.post(`/projects/${projectId}/assign`, {
            engineer_id: parseInt(engineerId, 10)
        });
        return response.data;
    }
    catch (err){
        console.error(err.response || err.message || err)
        throw new Error(err.response?.data?.error || err.message)
    }
}

export const exportDefectsCSV = async (projectId) => {
    try {
        const response = await api.axiosInstance.get(`/projects/${projectId}/export`, {
            responseType: 'blob'
        });

        const url = window.URL.createObjectURL(new Blob([response.data], { type: 'text/csv' }));
        const link = document.createElement('a');
        link.href = url;
        link.setAttribute('download', `project_${projectId}_defects_report.csv`);
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
        window.URL.revokeObjectURL(url);

        return { message: 'CSV exported successfully' };
    } catch (err) {
        console.error('Error in exportDefectsCSV:', err.response || err.message || err);
        throw new Error(err.response?.data?.error || 'Failed to export CSV');
    }
};