"use client"

import {ReactNode} from 'react';
import {useModal} from "@/states/modalStateHandler";

export default function Modal({children, width, height}: {children: ReactNode, width: string, height: string}) {
    const toggleModal = useModal((state) => state.toggleModal)
    const className = `bg-neutral-700 p-2 rounded-xl w-[${width}] h-[${height}] overflow-auto px-20 py-15`

    return (
        <div className="fixed backdrop-blur-sm inset-0 bg-opacity-25 bg-black flex justify-center items-center">
            <div className="flex flex-col">
                <button onClick={() => toggleModal(false)} className="text-white text-xl place-self-end">X</button>
                <div data-testid="modal" className={className}>{children}</div>
            </div>
        </div>
    )
}