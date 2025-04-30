"use client"

import Navbar from "@/components/navbar";
import {redirect} from "next/navigation";
import {LinkTo} from "@/components/pageEnums";
import {useEffect} from "react";


export default function Layout({children}: Readonly<{children: React.ReactNode}>) {
    useEffect(() => {
        const value = localStorage.getItem("sessionToken") || ""
        if (!value) redirect(LinkTo.LOGIN)
    })

    return (
        <div className="flex h-screen flex-col md:flex-row md:overflow-hidden">
            <div className="w-full flex-none md:w-64">
                <Navbar/>
            </div>
            <div className="flex-grow p-6 md:overflow-y-auto md:p-12">{children}</div>
        </div>
    )
}