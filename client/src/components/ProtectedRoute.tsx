import React from "react";
import { Navigate } from "react-router";
import { useAuth } from "../hooks/useAuth";
import { LoaderCircle } from "lucide-react";

interface ProtectedRouteProps {
    children: React.ReactNode;
}

const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ children }) => {
    const { isAuthenticated, isLoading } = useAuth();

    if (isLoading) {
        return (
            <div className="flex items-center justify-center min-h-screen">
                <div className="text-center flex items-center flex-col">
                    <LoaderCircle className="animate-spin w-16 h-16" />
                    <p className="mt-4 text-gray-600 text-2xl">Loading...</p>
                </div>
            </div>
        );
    }

    if (!isAuthenticated) {
        return <Navigate to="/login" />;
    }

    return <>{children}</>;
};

export default ProtectedRoute;
