"use client"

import {useSearchParams} from "next/navigation";
import useSWR from "swr";

export default function ListEmailsTable() {
    const searchParams = useSearchParams()
    const id = searchParams.get("id")

    const fetcher = (...args: Parameters<typeof fetch>) =>
        fetch(...args).then((res) => res.json())

    const {data, error} = useSWR(`https://jsonplaceholder.typicode.com/posts/${id}`, fetcher)

    if (error) return <div>Error</div>
    if (!data) {
        return <div>Loading...</div>
    } else if (!data.id) {
        return <div>Error</div>
    }




    // const [oldId, setId] = useState("")


    return (
        <>
            <div className="grid grid-cols-3">
                <h1 className="col-span-2 font-bold py-5 px-2 text-2xl">{data.title}</h1>
                <div className="py-5 px-2 flex flex-row justify-end">
                    <div className="px-3">
                    <a href={`/ui/mailingLists/list/manageMods.html?id=${id}`}
                       className="bg-cyan-600 text-white py-3 px-2 hover:bg-cyan-500 rounded-md">Manage Moderators</a>
                    </div>

                    <div>
                    <a href={`/ui/mailingLists/list/edit.html?id=${id}`}
                       className="bg-cyan-600 text-white py-3 px-2 hover:bg-cyan-500 rounded-md">Manage Senders</a>
                    </div>
                </div>
            </div>
            <div className="min-w-full table text-gray-900 shadow-[0_3px_10px_-1px_rgba(0,0,0,1)]">
                <div className="table-header-group text-left text-sm font-normal">
                    <div data-testid="table-head" className="table-row bg-neutral-800/45 text-neutral-300">
                        <div className="table-cell border-b border-black text-left px-4 py-5 font-bold sm:pl-6">
                            From
                        </div>
                        <div className="table-cell border-b border-black text-left px-4 py-5 font-bold sm:pl-6">
                            Title
                        </div>
                        <div className="table-cell border-b border-black text-left px-4 py-5 font-bold sm:pl-6">
                            Date
                        </div>
                    </div>
                </div>
                <div data-testid="table-body" className="table-row-group">
                    <a href={`/ui/mailingLists/list/email.html?listId=${1}?id=${1}`}
                          className="table-row shadow-inner text-neutral-300 hover:bg-neutral-600/75  hover:cursor-pointer">
                        <div className="table-cell border-black border-b py-3 text-sm">
                            <div className="whitespace-nowrap py-3 pl-6 pr-3 flex items-center gap-3">
                                me@example.com
                            </div>
                        </div>
                        <div className="table-cell border-black border-b py-3 text-sm">
                            <div className="whitespace-nowrap py-3 pl-6 pr-3 flex items-center gap-3">
                                An email
                            </div>
                        </div>
                        <div className="table-cell border-black border-b py-3 text-sm">
                            <div className="whitespace-nowrap py-3 pl-6 pr-3 flex items-center gap-3">
                                09/01/2025
                            </div>
                        </div>
                    </a>

                    <a href={`/ui/mailingLists/list/email.html?listId=${1}?id=${2}`}
                          className="table-row shadow-inner text-neutral-300 hover:bg-neutral-600/75  hover:cursor-pointer">
                        <div className="table-cell black border-black py-3 text-sm">
                            <div className="whitespace-nowrap py-3 pl-6 pr-3 flex items-center gap-3">
                                exampleentry2@email.com
                            </div>
                        </div>
                        <div className="table-cell black border-black py-3 text-sm">
                            <div className="whitespace-nowrap py-3 pl-6 pr-3 flex items-center gap-3">
                                Testing Testing...
                            </div>
                        </div>
                        <div className="table-cell border-black border-b py-3 text-sm">
                            <div className="whitespace-nowrap py-3 pl-6 pr-3 flex items-center gap-3">
                                09/01/2025
                            </div>
                        </div>
                    </a>
                </div>
            </div>
        </>
    )
}