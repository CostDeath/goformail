import { create } from "zustand"

interface ModalState {
    toggled: boolean;
    toggleModal: (toggleTo: boolean) => void;
}

export const useModal = create<ModalState>()((set) => ({
    toggled: false,
    toggleModal: (toggleTo) => set((state) => ({ toggled: toggleTo }))
}))