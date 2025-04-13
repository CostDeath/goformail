export default function MailingListTable({query, currentPage, api}: {
    query?: string;
    currentPage?: number;
    api: string;
}) {
    // having api be nullable for now, TODO: make it not nullable once api is ready
    // let data = await fetch(api);
    // let retriedData = data.json();
    console.log(query);
    console.log(currentPage);
    console.log(api);

    return (
        <div className="min-w-full table text-gray-900 shadow-[0_3px_10px_-1px_rgba(0,0,0,1)]">
            <div className="table-header-group text-left text-sm font-normal">
                <div data-testid="table-head" className="table-row bg-neutral-800/45 text-neutral-300">
                    <div className="table-cell border-b border-black text-left px-4 py-5 font-bold sm:pl-6">
                        Society
                    </div>
                    <div className="table-cell border-b border-black text-left px-4 py-5 font-bold sm:pl-6">
                        Mailing List Email
                    </div>
                </div>
            </div>
            <div data-testid="table-body" className="table-row-group">
                <div className="table-row shadow-inner text-neutral-300 hover:bg-neutral-600/75 hover:cursor-pointer">
                    <div className="table-cell border-black border-b py-3 text-sm">
                        <div className="whitespace-nowrap py-3 pl-6 pr-3 flex items-center gap-3">
                            Example Entry Society
                        </div>
                    </div>
                    <div className="table-cell border-black border-b py-3 text-sm">
                        <div className="whitespace-nowrap py-3 pl-6 pr-3 flex items-center gap-3">
                            exampleentry@email.com
                        </div>
                    </div>
                </div>

                <div className="table-row shadow-inner text-neutral-300 hover:bg-neutral-600/75 hover:cursor-pointer">
                    <div className="table-cell border-black border-b py-3 text-sm">
                        <div className="whitespace-nowrap py-3 pl-6 pr-3 flex items-center gap-3">
                            Example Entry 2 Society
                        </div>
                    </div>
                    <div className="table-cell border-black border-b py-3 text-sm">
                        <div className="whitespace-nowrap py-3 pl-6 pr-3 flex items-center gap-3">
                            exampleentry2@email.com
                        </div>
                    </div>
                </div>
            </div>
        </div>
    )
}