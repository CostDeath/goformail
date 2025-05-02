"use client"

import {useSearchParams} from "next/navigation";
import useSWR from "swr";
import {useEffect, useState} from "react";
import {getSessionToken} from "@/components/sessionToken";
import {api} from "@/components/api";
import {Email} from "@/models/email";

export default function ListEmailsTable() {
    const searchParams = useSearchParams()
    const id = searchParams.get("id")
    const [listName, setListName] = useState("")
    const [baseUrl, setBaseUrl] = useState("")
    const [offset, setOffset] = useState(0);
    const [sessionToken, setSessionToken] = useState<string | null>()

    const fetcher = async (url: string) => {
        let response = await fetch(`${baseUrl}${api.list}?id=${id}`, {
            method: "GET",
            headers: {
                "Authorization": `Bearer ${sessionToken}`
            }
        })

        if (response.ok) {
            const data = await response.json();
            setListName(data.data.name)
        }

        response = await fetch(url, {
            method: "POST",
            body: JSON.stringify({
                list: Number(id),
                archived: true,
                offset: offset
            }),
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${sessionToken}`
            }
        })
        if (response.ok) {
            return await response.json()
        } else {
            return await response.text()
        }
    }

    useEffect(() => {
        const url = `${window.location.origin}/api`
        setBaseUrl(url)
        setSessionToken(getSessionToken())
    }, [])


    const {data, error} = useSWR((baseUrl && sessionToken) ? `${baseUrl}${api.emails}` : null, fetcher)

    if (error) return <div>Error loading emails list</div>
    if (!data) {
        return <div>Loading...</div>
    } else if (data.message !== "Successfully fetched emails!") {
        return <div>Error loading emails list</div>
    }

    const emails: Email[] = data.data.emails

    return (
        <>
            <div className="grid grid-cols-3">
                <h1 className="col-span-2 font-bold py-5 px-2 text-2xl">{listName}</h1>
                <div className="py-5 px-2 flex flex-row justify-end">
                    <div className="px-3">
                        <a href={`/ui/mailingLists/list/manageMods.html?id=${id}`}
                           className="bg-cyan-600 text-white py-3 px-2 hover:bg-cyan-500 rounded-md">Manage
                            Moderators</a>
                    </div>

                    <div>
                        <a href={`/ui/mailingLists/list/edit.html?id=${id}`}
                           className="bg-cyan-600 text-white py-3 px-2 hover:bg-cyan-500 rounded-md">Manage List</a>
                    </div>
                </div>
            </div>

            <div className="overflow-auto max-h-[50vh] shadow-[0_3px_10px_-1px_rgba(0,0,0,1)]">
                <div className="min-w-full table text-gray-900">
                    <div className="table-header-group text-left text-sm font-normal">
                        <div data-testid="table-head-emails" className="table-row bg-neutral-800/45 text-neutral-300">
                            <div className="table-cell border-b border-black text-left px-4 py-5 font-bold sm:pl-6">
                                From
                            </div>
                            <div className="table-cell border-b border-black text-left px-4 py-5 font-bold sm:pl-6">
                                Date
                            </div>
                        </div>
                    </div>
                    <div data-testid="table-body-emails" className="table-row-group">

                        {emails.length > 0 && (
                            <>
                                {emails.map((email) => (
                                    <a key={email.id}
                                       href={`/ui/mailingLists/list/email.html?listId=${1}&id=${email.id}`}
                                       className="table-row shadow-inner text-neutral-300 hover:bg-neutral-600/75  hover:cursor-pointer">
                                        <div className="table-cell border-black border-b py-3 text-sm">
                                            <div className="whitespace-nowrap py-3 pl-6 pr-3 flex items-center gap-3">
                                                {email.sender}
                                            </div>
                                        </div>
                                        <div className="table-cell border-black border-b py-3 text-sm">
                                            <div className="whitespace-nowrap py-3 pl-6 pr-3 flex items-center gap-3">
                                                {email.received_at.getFullYear()}-{email.received_at.getMonth() + 1}-{email.received_at.getDate()}.
                                            </div>
                                        </div>
                                    </a>
                                ))}
                            </>
                        )}

                        {emails.length < 1 && (
                            <>
                                <div
                                    className="table-row shadow-inner text-neutral-300">
                                    <div className="table-cell border-black border-b py-3 text-sm">
                                        <div className="whitespace-nowrap py-3 pl-6 pr-3 flex items-center gap-3">
                                            No Data to show
                                        </div>
                                    </div>
                                </div>
                            </>
                        )}
                    </div>
                </div>
            </div>

            <div data-testid="pagination" className="flex justify-center py-5">
                {data.data.offset >= 2000 && (
                    <div onClick={() => setOffset(data.data.offset - 2000)} className="hover:cursor-pointer px-[3]">
                        <div className="flex h-10 w-10 items-center justify-center text-sm border border-neutral-700/25 rounded-md hover:bg-neutral-600">
                            {"<"}
                        </div>
                    </div>
                )}

                {data.data.offset > 0 && (
                    <div onClick={() => setOffset(data.data.offset)} className="hover:cursor-pointer px-[3]">
                        <div className="flex h-10 w-10 items-center justify-center text-sm border border-neutral-700/25 rounded-md hover:bg-neutral-600">
                            {">"}
                        </div>
                    </div>
                )}
            </div>

            </>
            )
            }