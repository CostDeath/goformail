"use client"

import {useState} from "react";
import {permissionsList} from "@/components/permissions";

export default function AddUserForm() {
    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")
    const [permissions, setPermissions] = useState(permissionsList)

    const handleChange = (index: number) => {
        const values = [...permissions]
        values[index].value = !values[index].value
        setPermissions(values)
    }

    const placeholder = () => {
        console.log(email)
        console.log(password)
        console.log(permissions)
    }

    return (
        <>
            <h1 className="text-2xl font-bold py-5">Create a user</h1>
            <hr/>
            <div className="grid grid-cols-2 py-10">
                <label className="text-xl p-1">Email Address</label>
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
                        placeholder:text-neutral-500
                        "
                    id="email"
                    type="email"
                    name="email"
                    value={email}
                    placeholder="Email Address"
                    onChange={e => setEmail(e.target.value)}
                    required
                />
            </div>

            <div className="grid grid-cols-2">
                <label className="text-xl p-1">Password</label>
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
                        placeholder:text-neutral-500
                        "
                    id="password"
                    type="password"
                    name="password"
                    value={password}
                    placeholder="Password"
                    onChange={e => setPassword(e.target.value)}
                    required
                />
            </div>
            <br/>
            <hr/>
            <br/>
            <h1 className="text-2xl font-bold underline">Permissions</h1>

            <div className="py-5">
                {permissions.map((permission, index) => (
                    <div className="grid grid-cols-3 px-2 py-3" key={index}>
                        <label className="text-xl">{permission.label}</label>
                        <input id="permission" type="checkbox" value="" onClick={() => handleChange(index)} />
                    </div>
                ))}
            </div>

            <div className="flex flex-row justify-end px-5">
                <button className="bg-green-600/75 hover:bg-green-600 px-2 py-1 rounded-md"
                        onClick={placeholder}>
                    Submit
                </button>
            </div>
        </>
    )
}