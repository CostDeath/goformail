"use client"

import {ReactNode} from 'react';
import {useModal} from "@/states/modalStateHandler";

export default function Modal({children}: {children: ReactNode}) {
    const toggleModal = useModal((state) => state.toggleModal)

    return (
        <div className="fixed backdrop-blur-sm inset-0 bg-opacity-25 bg-black flex justify-center items-center">
            <div className="flex flex-col">
                <button onClick={() => toggleModal(false)} className="text-white text-xl place-self-end">X</button>
                <div data-testid="modal" className="bg-neutral-700 p-2 rounded-xl w-[150vh] h-[90vh] overflow-auto px-20 py-15">{children}</div>
            </div>
        </div>
    )
}