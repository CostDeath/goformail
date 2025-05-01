import {useEffect, useState} from "react";
import {getSessionToken} from "@/components/sessionToken";
import useSWR from "swr";
import {api} from "@/components/api";
import {useSearchParams} from "next/navigation";


export default function RecipientsTable() {
    const param = useSearchParams()
    const listId = param.get("id")

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

    const {data, error} = useSWR((baseUrl) ? `${baseUrl}${api.list}?id=${listId}` : null, fetcher)

    if (error) return <div>Error</div>
    if (!data) {
        return <div>Loading</div>
    } else if (data.message !== "Successfully fetched list!") return <div>Error</div>

    const list = data.data

    return (
        <>
            <div className="py-5 px-2 flex flex-row justify-end">
                <a href={`/ui/mailingLists/list/manageRecipients.html?id=${listId}`}
                   className="bg-cyan-600 text-white py-3 px-2 hover:bg-cyan-500 rounded-md">Manage Recipients</a>
            </div>
            <div className="overflow-auto max-h-[50vh] shadow-[0_3px_10px_-1px_rgba(0,0,0,1)]">
                <div className="min-w-full table text-gray-900">
                    <div className="table-header-group text-left text-sm font-normal">
                        <div data-testid="table-head" className="table-row bg-neutral-800/45 text-neutral-300">
                            <div className="table-cell border-b border-black text-left px-4 py-5 font-bold sm:pl-6">
                                Recipients
                            </div>
                        </div>
                    </div>
                    <div data-testid="table-body" className="table-row-group">

                        {list.recipients && (
                            <>
                                {list.recipients.map((recipient: string) => (
                                    <div key={recipient}
                                         className="table-row shadow-inner text-neutral-300">
                                        <div className="table-cell black border-black py-3 text-sm">
                                            <div className="whitespace-nowrap py-3 pl-6 pr-3 flex items-center gap-3">
                                                {recipient}
                                            </div>
                                        </div>
                                    </div>
                                ))}
                            </>
                        )}

                        {!list.recipients && (
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
        </>
    )
}