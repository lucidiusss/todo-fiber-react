import Tasks from "./components/Tasks";
import { Input } from "@/components/ui/input";
import "./index.css";
import { useEffect, useState } from "react";
import axios from "axios";
import { Button } from "./components/ui/button";
import { Plus } from "lucide-react";
import toast from "react-hot-toast";
import { Spinner } from "./components/ui/spinner";
import Header from "./components/Header";
import { api } from "./api/client";
import { useAuth } from "./hooks/useAuth";

export interface TaskInterface {
    id: number;
    title: string;
    completed: boolean;
    created_at: Date;
    updated_at: Date;
    deleted_at: Date;
}

function App() {
    const [tasks, setTasks] = useState<TaskInterface[]>([]);
    const [inputValue, setInputValue] = useState<string>("");
    const [isAdding, setIsAdding] = useState<boolean>(false);
    const [isLoading, setIsLoading] = useState<boolean>(true);

    const { user } = useAuth();

    useEffect(() => {
        getTasks();
    }, []);

    const getTasks = async () => {
        try {
            const { data } = await api.get("tasks");
            setTasks(data.data);
            setIsLoading(false);
        } catch (error) {
            toast.dismiss();
            if (axios.isAxiosError(error)) {
                const errorMessage = error.response?.data?.error;
                console.log(errorMessage);
                toast(`❌ ${errorMessage}`);
            }
        }
    };

    const createTask = async (title: string) => {
        try {
            setIsAdding(true);
            const newTitle = title.trim();
            const { data } = await api.post("/tasks", {
                title: newTitle,
                user_id: user?.id,
            });
            setTasks([...tasks, data.data]);
            toast("✅ New task is created!");
        } catch (error) {
            toast.dismiss();
            if (axios.isAxiosError(error)) {
                const errorMessage = error.response?.data?.error;
                console.log(errorMessage);
                toast(`❌ ${errorMessage}`);
            }
        } finally {
            setInputValue("");
            setIsAdding(false);
        }
    };

    return (
        <div className="flex h-screen items-center justify-center">
            <Header />
            <div className="max-w-7xl w-full md:my-10 md:w-2/3 h-3/4 md:min-h-1/2 flex flex-col items-center gap-10 p-3 md:p-6 rounded-xl shadow-xs bg-slate-50">
                <div>
                    <h1 className="text-2xl lg:text-3xl sm:text-3xl font-bold">
                        all tasks
                    </h1>
                </div>
                <div className="w-full md:w-3/4 flex flex-row items-center gap-5">
                    <Input
                        disabled={isAdding}
                        onKeyDown={(e) =>
                            e.key === "Enter" ? createTask(inputValue) : ""
                        }
                        value={inputValue}
                        onChange={(e) => setInputValue(e.target.value)}
                        placeholder="add new task"
                        className="placeholder:text-gray-300 sm:placeholder:text-2xl text-[18px] lg:text-[24px] lg:placeholder:text-[24px] md:placeholder:text-[20px] md:text-[20px] sm:py-3 sm:px-4 bg-gray-100 rounded-md"
                    />
                    {inputValue.length > 0 ? (
                        <Button
                            disabled={isAdding}
                            size="icon"
                            onClick={() => createTask(inputValue.trim())}
                            className="bg-gray-100 shadow border group  active:bg-green-500 hover:bg-gray-200"
                        >
                            {isAdding ? (
                                <Spinner className="text-black group-active:text-white" />
                            ) : (
                                <Plus
                                    size={18}
                                    className="text-black group-active:text-white"
                                />
                            )}
                        </Button>
                    ) : (
                        <></>
                    )}
                </div>
                <Tasks
                    isLoading={isLoading}
                    tasks={tasks}
                    setTasks={setTasks}
                />
            </div>
        </div>
    );
}

export default App;
