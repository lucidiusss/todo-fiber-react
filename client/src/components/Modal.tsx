import type { Dispatch, FC, SetStateAction } from "react";

interface ModalProps {
    modal: boolean;
    setModal: Dispatch<SetStateAction<boolean>>;
}

const Modal: FC<ModalProps> = ({ setModal, modal }) => {
    return (
        <>
            <div className="max-w-7xl z-10 w-full md:my-10 md:w-2/3 h-3/4 md:min-h-1/2 flex flex-col items-center gap-10 p-3 md:p-6 rounded-xl shadow-xs bg-slate-50"></div>
            <div className="absolute top-0 left-0 w-full z-0 flex items-center justify-center h-full bg-black/60" />
        </>
    );
};

export default Modal;
