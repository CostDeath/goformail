"use client"

import Link from "next/link";

const placeholder = () => {
    console.log("login placeholder");
}

export default function LoginForm() {
    // TODO: need to add logic where it validates email and password before changing page to mailing list page
    return (
        <form className="space-y-3" action={placeholder}>
            <h1 className="text-xl">
                Log In
            </h1>
            <div className="w-full">
                <div>
                    <label
                        className="mb-3 mt-5 block text-xs font-medium text-gray-900"
                        htmlFor="email"
                        >
                        Email
                    </label>
                    <div className="relative">
                        <input
                            className="peer block w-full rounded-md border border-gray-200 py-[9px] pl-10 text-sm outline-2 placeholder:text-gray-500"
                            id="email"
                            type="email"
                            name="email"
                            placeholder="Enter your email"
                            required
                            />
                        {/*Icon here potentially? */}
                    </div>
                </div>
                <div className="mt-4">
                    <label
                        className="mb-3 mt-5 block text-xs font-medium text-gray-900"
                        htmlFor="password"
                        >
                        Password
                    </label>
                    <div className="relative">
                        <input
                            className="peer block w-full rounded-mb border border-gray-200 py-[9px] pl-10 text-smm outline-2 placeholder:text-gray-500"
                            id="password"
                            type="password"
                            name="password"
                            placeholder="Enter your password"
                            required
                            />
                        {/*Icon here potentially? */}
                    </div>
                </div>
            </div>
            <button className="bg-blue-500 hover:bg-blue-600 text-white rounded-xl py-2 px-6">
                <Link href="/mailingLists">Log in</Link>
            </button>
        </form>
    )
}