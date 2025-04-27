"use client"

import useSWR from "swr";

export default function LogsTable() {
    const fetcher = async(url: string) => {
        const response = await fetch(url);
        return await response.json();
    }

    const {data, error} = useSWR("https://jsonplaceholder.typicode.com/comments", fetcher)

    if (error) return <div>Error</div>;
    if (!data) {
        return <div>Loading</div>
    }

    return (
        <div className="overflow-auto h-[75vh] shadow-[0_3px_10px_-1px_rgba(0,0,0,1)]">
        <div className="min-w-full table text-gray-900">
            <div className="table-header-group text-left text-sm font-normal">
                <div data-testid="table-head" className="table-row bg-neutral-800/45 text-neutral-300">
                    <div className="table-cell border-b border-black text-left px-4 py-5 font-bold sm:pl-6">
                        Date
                    </div>
                    <div className="table-cell border-b border-black text-left px-4 py-5 font-bold sm:pl-6">
                        Action
                    </div>
                    <div className="table-cell border-b border-black text-left px-4 py-5 font-bold sm:pl-6">
                        User
                    </div>
                </div>
            </div>
            <div data-testid="table-body" className="table-row-group">
                {data.map((log: any) => (
                    <div key={log.id} className="table-row shadow-inner text-neutral-300 hover:bg-neutral-600/75">
                        <div className="table-cell border-black border-b py-3 text-sm">
                            <div className="whitespace-nowrap py-3 pl-6 pr-3 flex items-center gap-3">
                                {log.id}
                            </div>
                        </div>
                        <div className="table-cell border-black border-b py-3 text-sm">
                            <div className="whitespace-nowrap py-3 pl-6 pr-3 flex items-center gap-3">
                                {log.name}
                            </div>
                        </div>
                        <div className="table-cell border-black border-b py-3 text-sm">
                            <div className="whitespace-nowrap py-3 pl-6 pr-3 flex items-center gap-3">
                                {log.email}
                            </div>
                        </div>
                    </div>
                ))}
            </div>
        </div>
        </div>
    )
}