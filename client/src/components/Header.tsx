import { useAuth } from "@/hooks/useAuth";
import { Button } from "./ui/button";

const Header = () => {
    const { logout, user } = useAuth();
    return (
        <div className="absolute top-5 flex items-center px-4 w-full">
            <h1 className="text-center text-2xl select-none">
                welcome, <span className="font-bold">{user?.username}</span>!
            </h1>
            <Button
                onClick={logout}
                className="hover:underline ml-auto cursor-pointer p-0 text-2xl"
                variant="ghost"
            >
                log out
            </Button>
        </div>
    );
};

export default Header;
