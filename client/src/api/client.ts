import axios from "axios";

export const api = axios.create({
    baseURL: import.meta.env.VITE_API_URL,
    headers: {
        "Content-Type": "application/json",
    },
});

// Request interceptor to add token to EVERY request
api.interceptors.request.use(
    (config) => {
        const token = localStorage.getItem("token");

        if (token) {
            // Make sure it's formatted exactly as "Bearer <token>"
            config.headers.Authorization = `Bearer ${token}`;
        } else {
            console.warn("No token found for request:", config.url);
        }

        return config;
    },
    (error) => {
        return Promise.reject(error);
    },
);

// Response interceptor for debugging
/* api.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error.response?.status === 401) {
            console.error("Auth error:", error.response.data);

            // If token is invalid/expired, redirect to login
            if (
                error.response.data?.error?.includes("Bearer") ||
                error.response.data?.error?.includes("token")
            ) {
                localStorage.removeItem("token");
                window.location.href = "/login";
            }
        }
        return Promise.reject(error);
    },
);
*/
