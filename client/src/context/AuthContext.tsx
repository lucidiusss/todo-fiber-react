import React, { useState, useEffect } from "react";
import type { ReactNode } from "react";
import { api } from "../api/client";
import { AuthContext } from "./auth.context";
import axios from "axios";

interface User {
    id: string;
    username: string;
}

export const AuthProvider: React.FC<{ children: ReactNode }> = ({
    children,
}) => {
    const [user, setUser] = useState<User | null>(null);
    const [token, setToken] = useState<string | null>(
        localStorage.getItem("token"),
    );
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        const loadUser = async () => {
            if (token) {
                try {
                    const response = await api.get("/profile");
                    setUser(response.data);
                } catch (error) {
                    console.error("Failed to load user:", error);
                    localStorage.removeItem("token");
                    setToken(null);
                }
            }
            setIsLoading(false);
        };

        loadUser();
    }, [token]);

    const login = async (username: string, password: string) => {
        try {
            const response = await api.post("/auth/login", {
                username,
                password,
            });
            const { token, user } = response.data;

            localStorage.setItem("token", token);
            setToken(token);
            setUser(user);
        } catch (error) {
            if (axios.isAxiosError(error)) {
                // Re-throw with the error message from backend
                throw new Error(error.response?.data?.message);
            }
            throw error; // Re-throw any other errors
        }
    };

    const register = async (username: string, password: string) => {
        try {
            const response = await api.post("/auth/register", {
                username,
                password,
            });
            const { token, user } = response.data;

            localStorage.setItem("token", token);
            setToken(token);
            setUser(user);
        } catch (error) {
            if (error instanceof Error) {
                throw new Error(error.message);
            }
            throw new Error("Registration failed");
        }
    };

    const logout = () => {
        localStorage.removeItem("token");
        setToken(null);
        setUser(null);
    };

    const value = {
        user,
        token,
        isLoading,
        login,
        register,
        logout,
        isAuthenticated: !!user,
    };

    return (
        <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
    );
};
export { AuthContext };
