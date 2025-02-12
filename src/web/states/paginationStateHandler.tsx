import { create } from "zustand"

interface PaginationState {
    pageNumber: number;
    changePage: (newNumber: number) => void;
    reset: () => void;
}

export const togglePagination = create<PaginationState>()((set) => ({
    pageNumber: 1,
    changePage: (newNumber) => set((state) => ({ pageNumber: newNumber })),
    reset: () => set((state) => ({pageNumber: 1}))
}))