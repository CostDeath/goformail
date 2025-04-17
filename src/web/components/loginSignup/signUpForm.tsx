"use client"

import StyledInput from "@/components/styledInput";

const placeholder = () => {
    console.log("sign up placeholder");
}

export default function SignUpForm() {
    return (
        <form className="space-y-3" action={placeholder}>
            <h1 className="text-xl">
                Sign Up
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
                        htmlFor="studentID"
                    >
                        Student ID
                    </label>
                    <div className="relative">
                        <StyledInput
                            id="studentID"
                            name="studentID"
                            placeholder="Enter your student ID"
                            required
                        />
                        {/*Icon here potentially? */}
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
                        {/*Icon here potentially? */}
                    </div>
                </div>
            </div>
            <div className="py-3">
                <button className="bg-cyan-600 hover:bg-cyan-700 text-gray-100 rounded-xl py-2 px-6">Sign up</button>
            </div>
        </form>
    )
}