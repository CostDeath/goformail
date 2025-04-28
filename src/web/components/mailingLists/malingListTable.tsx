"use client"

import Link from "next/link";
import useSWR from "swr";
import {MailingLists} from "@/models/list";

export default function MailingListTable({api}: {
    api: string;
}) {
    const fetcher = async(url: string) => {
        const response = await fetch(url)
        return await response.json()
    }

    const {data, error} = useSWR(api, fetcher)

    if (error) {
        return <div>Error</div>
    }
    if (!data) {
        return <div>Loading</div>
    } else if (data.message !== "Successfully fetched lists!") {
        console.log(data)
        return <div>Error</div>
    }
    // console.log(data)

    return (
        <div className="overflow-auto max-h-[75vh] shadow-[0_3px_10px_-1px_rgba(0,0,0,1)]">
        <div className="min-w-full table text-gray-900">
            <div className="table-header-group text-left text-sm font-normal">
                <div data-testid="table-head" className="table-row bg-neutral-800/45 text-neutral-300">
                    <div className="table-cell border-b border-black text-left px-4 py-5 font-bold sm:pl-6">
                        Mailing List Email
                    </div>
                </div>
            </div>
            <div data-testid="table-body" className="table-row-group">
                {data.data && (
                        <>
                            {data.data.map((list: MailingLists) => (
                                <Link key={list.id} href={`/mailingLists/list.html?id=${list.id}`} className="table-row shadow-inner text-neutral-300 hover:bg-neutral-600/75 hover:cursor-pointer">
                                    <div className="table-cell border-black border-b py-3 text-sm">
                                        <div className="whitespace-nowrap py-3 pl-6 pr-3 flex items-center gap-3">
                                            {list.list.name}
                                        </div>
                                    </div>
                                </Link>
                        ))}
                        </>
                )
                }
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