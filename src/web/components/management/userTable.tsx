"use client"

import {useEffect, useState} from "react";
import useSWR from "swr";
import {api} from "@/components/api";
import {User} from "@/models/user";
import {getSessionToken} from "@/components/sessionToken";

export default function UserTable() {
    const [sessionToken, setSessionToken] = useState<string | null>();
    const fetcher = async(url: string) => {
        const response = await fetch(url, {
            method: "GET",
            headers: {
                "Authorization": `Bearer ${sessionToken}`
            }
        })
        return await response.json()
    }

    const [baseUrl, setBaseUrl] = useState("")

    useEffect(() => {
        const url = `${window.location.origin}/api`
        setBaseUrl(url)
        setSessionToken(getSessionToken())
    }, [])

    const {data, error} = useSWR((baseUrl) ? `${baseUrl}${api.users}` : null, fetcher)

    if (error) return <div>Error</div>
    if (!data) {
        return <div>Loading</div>
    } else if (data.message !== "Successfully fetched users!") {
        return <div>Error</div>
    }

    return(
        <div className="overflow-auto max-h-[75vh] shadow-[0_3px_10px_-1px_rgba(0,0,0,1)]">
            <div className="min-w-full table text-gray-900">
                <div className="table-header-group text-left text-sm font-normal">
                    <div data-testid="table-head" className="table-row bg-neutral-800/45 text-neutral-300">
                        <div className="table-cell border-b border-black text-left px-4 py-5 font-bold sm:pl-6">
                            User
                        </div>
                    </div>
                </div>
                <div data-testid="table-body" className="table-row-group">

                    {data.data && (
                        <>
                            {data.data.map((user: User) => (
                                // TODO: permission logic of being able to view users' perms
                                <a key={user.id} href={`/ui/management/edit.html?id=${user.id}`}
                                   className="table-row shadow-inner text-neutral-300 hover:bg-neutral-600/75  hover:cursor-pointer">
                                    <div className="table-cell black border-black py-3 text-sm">
                                        <div className="whitespace-nowrap py-3 pl-6 pr-3 flex items-center gap-3">
                                            {user.email}
                                        </div>
                                    </div>
                                </a>
                            ))}
                        </>
                    )}

                    {!data.data && (
                        <div className="table-row shadow-inner text-neutral-300">
                            <div className="table-cell border-black border-b py-3 text-sm">
                                <div className="whitespace-nowrap py-3 pl-6 pr-3 flex items-center gap-3">
                                    No Data to Show
                                </div>
                            </div>
                        </div>
                    )}
                </div>
            </div>
        </div>
    )
}