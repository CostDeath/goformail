"use client"


import {useSearchParams} from "next/navigation";
import useSWR from "swr";
import {useState} from "react";
import {permissionsList} from "@/components/permissions";
import DeleteUser from "@/components/editUser/deleteUser";

export default function EditUserForm() {
    // TODO: Replace "permissionsList" with the appropriate data when doing fetching ticket
    const [permissions, setPermissions] = useState(permissionsList)
    const search = useSearchParams()
    const userId = search.get("id")


    const fetcher = async(url: string) => {
        const response = await fetch(url)
        return await response.json()
        // setPermissions(data.permissions)   something like that
    }

    const {data, error} = useSWR(`https://jsonplaceholder.typicode.com/posts/${userId}`, fetcher)

    if (error) return <div>Error</div>
    if (!data) {
        return <div>Loading</div>
    } else if (!data.id) {
        return <div>Error</div>
    }



    const handleChange = (index: number) => {
        const values = [...permissions]
        values[index].value = !values[index].value
        setPermissions(values)
    }

    const placeholder = () => {
        console.log(data.email) // assuming data given, email is name of key
        console.log(permissions)
    }

    return (
        <>
            <h1 className="text-2xl font-bold py-5">Edit User</h1>
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
                    value={data.email}
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
                        <label className="text-xl">{permission.label}</label>
                        <input id="permission" type="checkbox" value="" onClick={() => handleChange(index)} />
                    </div>
                ))}
            </div>

            <div className="flex flex-row justify-end px-5">
                <div className="grid grid-cols-2">
                    <div>
                        <button className="bg-green-600/75 hover:bg-green-600 px-3 py-2 rounded-md"
                            onClick={placeholder}>
                        Edit User
                        </button>
                    </div>
                    <DeleteUser />
                </div>
            </div>
        </>
    )
}