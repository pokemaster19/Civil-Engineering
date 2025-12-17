import api from "./axiosInstance";

export const uploadAttachment = async (formData) => {
  try {
    const response = await api.axiosInstance.post(`/attachments/`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return response.data;
  } catch (err) {
    console.error('Ошибка при загрузке вложения:', err);
    throw err; 
  }
};