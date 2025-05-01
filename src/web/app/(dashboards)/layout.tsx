"use client"

import Navbar from "@/components/navbar";
import {redirect} from "next/navigation";
import {LinkTo} from "@/components/pageEnums";
import {useEffect} from "react";
import {api} from "@/components/api";
import {handleInvalidSessionToken} from "@/components/sessionToken";


export default function Layout({children}: Readonly<{children: React.ReactNode}>) {

    const validateToken = async (url: string, token: string) => {
        const response = await fetch(url, {
            method: "POST",
            headers: {
                "Authorization": `Bearer ${token}`,
            }
        })
        if (!response.ok) {
            alert("Your session has expired, please log back in")
            handleInvalidSessionToken()
        }
    }

    useEffect(() => {
        const value = localStorage.getItem("sessionToken") || ""
        if (!value) redirect(LinkTo.LOGIN)
        const url = `${window.location.origin}/api${api.tokenValidation}` || ""
        validateToken(url, value)
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