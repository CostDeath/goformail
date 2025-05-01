"use client"

import {useState} from "react";
import {permissionsList} from "@/components/permissions";
import {api} from "@/components/api";
import validateEmail from "@/components/validateEmails";
import {redirect} from "next/navigation";
import {LinkTo} from "@/components/pageEnums";
import {getSessionToken} from "@/components/sessionToken";

export default function AddUserForm() {
    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")
    const [permissions, setPermissions] = useState(permissionsList)

    const handleChange = (index: number) => {
        const values = [...permissions]
        values[index].value = !values[index].value
        setPermissions(values)
    }

    const createUser = async () => {
        const url = `${window.location.origin}/api${api.user}`
        const sessionToken = getSessionToken()
        if (!validateEmail(email)) {
            alert(`${email} is not a valid email!`)
            return
        }
        const perms: string[] = []
        for (let i = 0; i < permissions.length; i++) {
            if (permissions[i].value) {
                perms.push(permissions[i].id)
            }
        }

        const response = await fetch(url, {
            method: "POST",
            body: JSON.stringify({
                Email: email,
                Password: password,
                Permissions: perms
            }),
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${sessionToken}`
            }
        })

        if (response.ok) {
            const result = await response.json()
            alert(result.message)
            redirect(LinkTo.MANAGEMENT)
        } else {
            const result = await response.text()
            alert(result)
        }
    }

    return (
        <>
            <h1 className="text-2xl font-bold py-5">Create a user</h1>
            <hr/>
            <div className="grid grid-cols-2 py-10">
                <label htmlFor="email" className="text-xl p-1">Email Address</label>
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
                <label htmlFor="password" className="text-xl p-1">Password</label>
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
                        <label htmlFor={permission.id} className="text-xl">{permission.label}</label>
                        <input id={permission.id} name={permission.id} type="checkbox" value="" onClick={() => handleChange(index)} />
                    </div>
                ))}
            </div>

            <div className="flex flex-row justify-end px-5">
                <button className="bg-green-600/75 hover:bg-green-600 px-2 py-1 rounded-md"
                        onClick={createUser}>
                    Create User
                </button>
            </div>
        </>
    )
}