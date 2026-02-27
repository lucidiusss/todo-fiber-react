import axios from "axios";

// Configure axios defaults
export const api = axios.create({
    baseURL: import.meta.env.VITE_API_URL,
});

// Add token to requests if it exists
api.interceptors.request.use((config) => {
    const token = localStorage.getItem("token");
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});
