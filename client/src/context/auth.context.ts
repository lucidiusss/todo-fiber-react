import { createContext } from "react";

interface User {
    id: string;
    username: string;
}

export interface AuthContextType {
    user: User | null;
    token: string | null;
    isLoading: boolean;
    login: (email: string, password: string) => Promise<void>;
    register: (email: string, password: string, name?: string) => Promise<void>;
    logout: () => void;
    isAuthenticated: boolean;
}

// Create and export the context separately
export const AuthContext = createContext<AuthContextType | undefined>(
    undefined,
);
