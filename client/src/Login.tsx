import { Field, FieldGroup, FieldLabel, FieldSet } from "@/components/ui/field";
import { Input } from "./components/ui/input";
import { Button } from "./components/ui/button";
import { Spinner } from "./components/ui/spinner";
import { useState } from "react";
import { NavLink } from "react-router";
import { EyeIcon, EyeOffIcon } from "lucide-react";
import toast from "react-hot-toast";
import axios from "axios";
import { useNavigate, Link, useLocation } from "react-router";
import { useAuth } from "./hooks/useAuth";

const Login = () => {
    const [isLogin, setIsLogin] = useState<boolean>(false);
    const [isPasswordVisible, setIsPasswordVisible] = useState(false);

    const [username, setUsername] = useState<string>("");
    const [password, setPassword] = useState<string>("");

    const navigate = useNavigate();
    const { login } = useAuth();

    const Login = async () => {
        // Validate inputs
        if (!password.trim() && !username.trim()) {
            toast.error("Username and password are required");
            return;
        }
        if (!username.trim()) {
            toast.error("Username is required");
            return;
        }
        if (!password.trim()) {
            toast.error("Password is required");
            return;
        }

        setIsLogin(true);
        const loadingToast = toast.loading("Logging in...");
        try {
            await login(username, password);

            toast.dismiss(loadingToast);
            toast.success("Login successful!");

            navigate("/");
        } catch (error) {
            toast.dismiss(loadingToast);

            if (axios.isAxiosError(error)) {
                // Handle Axios errors
                const errorMessage = error.response?.data?.message;
                toast.error(errorMessage);
                console.log("Error response:", error.response?.data);
            } else if (error instanceof Error) {
                toast.error(error.message);
            } else {
                toast.error("An unexpected error occurred");
            }
        } finally {
            setIsLogin(false);
        }
    };

    return (
        <div className="h-screen w-full flex items-center justify-center">
            <div className="bg-slate-50 p-10 rounded-xl shadow-xs flex flex-col justify-center items-center gap-5">
                <div>
                    <h1 className="font-bold text-2xl">
                        Sign in to manage your tasks
                    </h1>
                </div>
                <div>
                    <form onSubmit={(e) => e.preventDefault()}>
                        <FieldSet className="w-full max-w-xs">
                            <FieldGroup>
                                {/*username*/}
                                <Field>
                                    <FieldLabel
                                        className="text-[16px]"
                                        htmlFor="username"
                                    >
                                        Username
                                    </FieldLabel>
                                    <Input
                                        id="username"
                                        value={username}
                                        onChange={(e) =>
                                            setUsername(e.target.value)
                                        }
                                        onKeyDown={(e) =>
                                            e.key === "Enter" ? Login() : ""
                                        }
                                        type="text"
                                        placeholder="enter your username"
                                    />
                                </Field>
                                {/*password*/}
                                <Field>
                                    <FieldLabel
                                        className="text-[16px]"
                                        htmlFor="password"
                                    >
                                        Password
                                    </FieldLabel>

                                    <div className="relative">
                                        <Input
                                            id="password"
                                            value={password}
                                            onChange={(e) =>
                                                setPassword(e.target.value)
                                            }
                                            onKeyDown={(e) =>
                                                e.key === "Enter" ? Login() : ""
                                            }
                                            type={
                                                isPasswordVisible
                                                    ? "text"
                                                    : "password"
                                            }
                                            placeholder="••••••••••••••"
                                        />
                                        <Button
                                            variant="ghost"
                                            tabIndex={-1}
                                            size="icon"
                                            onMouseDown={() =>
                                                setIsPasswordVisible(
                                                    (prevState) => !prevState,
                                                )
                                            }
                                            className="cursor-pointer text-muted-foreground focus-visible:ring-ring/50 absolute inset-y-0 right-0 rounded-l-none hover:bg-transparent"
                                        >
                                            {isPasswordVisible ? (
                                                <EyeOffIcon size={16} />
                                            ) : (
                                                <EyeIcon size={16} />
                                            )}
                                            <span className="sr-only">
                                                {isPasswordVisible
                                                    ? "Hide password"
                                                    : "Show password"}
                                            </span>
                                        </Button>
                                    </div>
                                </Field>
                            </FieldGroup>
                        </FieldSet>
                        <Button
                            disabled={isLogin}
                            onClick={() => Login()}
                            className="w-full mt-10 cursor-pointer"
                        >
                            {isLogin ? "Signing in..." : "Sign in"}
                            {isLogin ? <Spinner /> : <></>}
                        </Button>
                    </form>

                    <div className="flex items-center justify-center mt-5">
                        <h1 className="text-black/50">
                            Dont have an account yet?
                            <NavLink
                                to="/register"
                                className="ml-2 text-black hover:underline"
                            >
                                Sign up
                            </NavLink>
                        </h1>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Login;
