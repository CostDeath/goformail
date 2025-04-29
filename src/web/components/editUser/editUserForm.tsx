"use client"


import {redirect, useSearchParams} from "next/navigation";
import useSWR from "swr";
import {useEffect, useState} from "react";
import {permissionsList} from "@/components/permissions";
import DeleteUser from "@/components/editUser/deleteUser";
import {api} from "@/components/api";
import {User} from "@/models/user";
import {LinkTo} from "@/components/pageEnums";

export default function EditUserForm() {
    const [permissions, setPermissions] = useState(permissionsList)
    const search = useSearchParams()
    const userId = search.get("id")

    const handleChange = (index: number) => {
        const values = [...permissions]
        values[index].value = !values[index].value
        setPermissions(values)
    }


    const fetcher = async(url: string) => {
        const response = await fetch(url)
        const data = await response.json()
        const dataPerms: string[] = data.data.permissions
        let perms = 0
        for (let i = 0; i < permissions.length; i++) {
            if (dataPerms[perms] === permissions[i].id) {
                handleChange(i)
                perms += 1
            }
        }
        return data
    }

    const [baseUrl, setBaseUrl] = useState("")

    useEffect(() => {
        const url = `${window.location.origin}/api`
        setBaseUrl(url)
    }, [])

    const {data, error} = useSWR((baseUrl) ? `${baseUrl}${api.user}?id=${userId}` : null, fetcher)

    if (error) return <div>Error</div>
    if (!data) {
        return <div>Loading</div>
    } else if (data.message !== "Successfully fetched user!") {
        return <div>Error</div>
    }

    const result: User = data.data


    const editUser = async () => {
        const url = `${window.location.origin}/api${api.user}?id=${userId}`
        const perms: string[] = []
        for (let i = 0; i < permissions.length; i++) {
            if (permissions[i].value) {
                perms.push(permissions[i].id)
            }
        }

        const response = await fetch(url, {
            method: "PATCH",
            body: JSON.stringify({
                Email: result.email,
                Permissions: perms
            }),
            headers: {
                "Content-Type": "application/json",
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
            <h1 className="text-2xl font-bold py-5">Edit User</h1>
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
                    value={result.email}
                    placeholder="Email Address"
                    disabled
                />
            </div>

            <br />
            <hr/>
            <br/>

            <h1 className="text-2xl font-bold underline">Permissions</h1>

            <div className="py-5">
                {permissions.map((permission, index) => (
                    <div className="grid grid-cols-3 px-2 py-3" key={index}>
                        <label htmlFor={permission.id} className="text-xl">{permission.label}</label>

                        {permission.value && (
                            <input checked id={permission.id} name={permission.id} type="checkbox" value=""
                                   onClick={() => handleChange(index)}/>
                        )}

                        {!permission.value && (
                            <input id={permission.id} name={permission.id} type="checkbox" value=""
                                   onClick={() => handleChange(index)}/>
                        )}

                    </div>
                ))}
            </div>

            <div className="flex flex-row justify-end px-5">
            <div className="grid grid-cols-2">
                <div>
                    <button className="bg-green-600/75 hover:bg-green-600 px-3 py-2 rounded-md"
                            onClick={editUser}>
                        Edit User
                    </button>
                </div>
                <DeleteUser id={result.id}/>
            </div>
            </div>
        </>
    )
}