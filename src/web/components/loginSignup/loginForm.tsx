"use client"

import { redirect } from "next/navigation";
import {LinkTo} from "@/components/pageEnums";
import StyledInput from "@/components/styledInput";

const placeholder = () => {
    console.log("login placeholder");
    redirect(LinkTo.MAILINGLISTS)
}

export default function LoginForm() {
    // TODO: need to add logic where it validates email and password before changing page to mailing list page
    return (
        <form className="space-y-3" action={placeholder}>
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
                        <StyledInput
                            id="email"
                            type="email"
                            name="email"
                            placeholder="Enter your email"
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
                        <StyledInput
                            id="password"
                            type="password"
                            name="password"
                            placeholder="Enter your password"
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