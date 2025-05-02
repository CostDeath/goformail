"use client"

import useSWR from "swr";
import {useEffect, useState} from "react";
import {getSessionToken} from "@/components/sessionToken";
import {api} from "@/components/api";
import {Email} from "@/models/email";


export default function ApprovalRequestsTable() {
    const [sessionToken, setSessionToken] = useState<string | null>();
    const [baseUrl, setBaseUrl] = useState<string>();
    const [offset, setOffset] = useState(0);

    useEffect(() => {
        const url = `${window.location.origin}/api`
        setBaseUrl(url);
        setSessionToken(getSessionToken());
    }, [])

    const fetcher = async (url: string) => {
        const response = await fetch(url, {
            method: "POST",
            body: JSON.stringify({
                offset: offset,
                pending_approval: true
            }),
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${sessionToken}`
            }
        });
        if (response.ok) {
            return await response.json();
        }
        return await response.text()
    }

    const {data, error} = useSWR((baseUrl && sessionToken) ? `${baseUrl}${api.emails}` : null, fetcher);

    if (error) return <div>Error</div>
    if (!data) return <div>Loading</div>
    if (data.message !== "Successfully fetched emails!") return <div>Error</div>

    const emails: Email[] = data.data.emails

    return (
        <>
            <div className="overflow-auto max-h-[75vh] shadow-[0_3px_10px_-1px_rgba(0,0,0,1)]">
                <div className="min-w-full table text-gray-900">
                    <div className="table-header-group text-left text-sm font-normal">
                        <div data-testid="table-head" className="table-row bg-neutral-800/45 text-neutral-300">
                            <div className="table-cell border-b border-black text-left px-4 py-5 font-bold sm:pl-6">
                                Mailing List Email {/*data.fact*/}
                            </div>
                            <div className="table-cell border-b border-black text-left px-4 py-5 font-bold sm:pl-6">
                                Date Sent for Approval
                            </div>
                        </div>
                    </div>
                    <div data-testid="table-body" className="table-row-group">


                        {emails.length > 0 && (
                            <>
                                {emails.map((email) => (
                                    <a key={email.id}
                                       href={`/ui/approvals/email.html?id=${email.id}`}
                                       className="table-row shadow-inner text-neutral-300 hover:bg-neutral-600/75  hover:cursor-pointer">
                                        <div className="table-cell border-black border-b py-3 text-sm">
                                            <div className="whitespace-nowrap py-3 pl-6 pr-3 flex items-center gap-3">
                                                {email.sender}
                                            </div>
                                        </div>
                                        <div className="table-cell border-black border-b py-3 text-sm">
                                            <div className="whitespace-nowrap py-3 pl-6 pr-3 flex items-center gap-3">
                                                {email.received_at}
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
                            <div
                                className="flex h-10 w-10 items-center justify-center text-sm border border-neutral-700/25 rounded-md hover:bg-neutral-600">
                                {"<"}
                            </div>
                        </div>
                    )}

                    {data.data.offset > 0 && (
                        <div onClick={() => setOffset(data.data.offset)} className="hover:cursor-pointer px-[3]">
                            <div
                                className="flex h-10 w-10 items-center justify-center text-sm border border-neutral-700/25 rounded-md hover:bg-neutral-600">
                                {">"}
                            </div>
                        </div>
                    )}
                </div>

            </>
            )
            }