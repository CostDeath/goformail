"use client"

import { redirect } from "next/navigation";
import {LinkTo} from "@/components/pageEnums";
import {useEffect, useState} from "react";
import {api} from "@/components/api";
import {getSessionToken} from "@/components/sessionToken";



export default function LoginForm() {
    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")

    const validateToken = async (url: string, token: string) => {
        const response = await fetch(url, {
            method: "POST",
            headers: {
                "Authorization": `Bearer ${token}`,
            }
        })
        if (response.ok) {
            redirect(LinkTo.MAILINGLISTS)
        }
    }

    useEffect(() => {
        const sessionToken = getSessionToken()
        const url = `${window.location.origin}/api${api.tokenValidation}`
        if (sessionToken) {
            validateToken(url, sessionToken)
        }
    })

    const login = async () => {
        const url = `${window.location.origin}/api${api.login}`
        const response = await fetch(url, {
            method: "POST",
            body: JSON.stringify({
                Email: email,
                Password: password,
            }),
            headers: {
                "Content-Type": "application/json"
            }
        })

        if (response.ok) {
            const result = await response.json()
            localStorage.setItem("sessionToken", result.data.token)
            redirect(LinkTo.MAILINGLISTS)
        } else {
            const result = await response.text()
            alert(result)
        }
    }

    return (
        <form className="space-y-3" action={login}>
            <h1 className="text-xl font-bold">
                Log In
            </h1>
            <div className="w-full">
                <div>
                    <label
                        className="mb-3 mt-5 block font-bold"
                        htmlFor="email"
                        >
                        Email
                    </label>
                    <div className="relative">
                        <input
                            className="
                            bg-neutral-700
                            peer
                            block
                            w-full
                            h-10
                            px-3
                            border
                          border-neutral-500
                            rounded-md
                            outline-2
                          placeholder:text-neutral-500"
                            id="email"
                            type="email"
                            name="email"
                            placeholder="Enter your email"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                            required
                            />
                    </div>
                </div>
                <div className="mt-4">
                    <label
                        className="mb-3 mt-5 block font-bold"
                        htmlFor="password"
                        >
                        Password
                    </label>
                    <div className="relative">
                        <input
                            className="
                            bg-neutral-700
                            peer
                            block
                            w-full
                            h-10
                            px-3
                            border
                            border-neutral-500
                            rounded-md
                            outline-2
                            placeholder:text-neutral-500"
                            id="password"
                            type="password"
                            name="password"
                            placeholder="Enter your password"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            required
                            />
                    </div>
                </div>
            </div>
            <div className="py-3">
                <button className="bg-cyan-600 hover:bg-cyan-700 text-gray-100 rounded-xl py-2 px-6">
                    Log in
                </button>
            </div>
        </form>
    )
}