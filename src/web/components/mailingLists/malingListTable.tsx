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
        <div className="min-w-full table text-gray-900 border-gray-400 border-[1px]">
            <div className="table-header-group text-left text-sm font-normal">
                <div data-testid="table-head" className="table-row">
                    <div className="table-cell border-b border-black text-left px-4 py-5 font-bold sm:pl-6">
                        Society
                    </div>
                    <div className="table-cell border-b border-black text-left px-4 py-5 font-bold sm:pl-6">
                        Mailing List Email
                    </div>
                </div>
            </div>
            <div data-testid="table-body" className="table-row-group">
                <div className="table-row shadow-inner hover:bg-gray-100 hover:cursor-pointer">
                    <div className="table-cell border-gray-200 border-b py-3 text-sm">
                        <div className="whitespace-nowrap py-3 pl-6 pr-3 flex items-center gap-3">
                            Example Entry Society
                        </div>
                    </div>
                    <div className="table-cell border-gray-200 border-b py-3 text-sm">
                        <div className="whitespace-nowrap py-3 pl-6 pr-3 flex items-center gap-3">
                            exampleentry@email.com
                        </div>
                    </div>
                </div>

                <div className="table-row shadow-inner hover:bg-gray-100 hover:cursor-pointer">
                    <div className="table-cell border-gray-200 border-b py-3 text-sm">
                        <div className="whitespace-nowrap py-3 pl-6 pr-3 flex items-center gap-3">
                            Example Entry 2 Society
                        </div>
                    </div>
                    <div className="table-cell border-gray-200 border-b py-3 text-sm">
                        <div className="whitespace-nowrap py-3 pl-6 pr-3 flex items-center gap-3">
                            exampleentry2@email.com
                        </div>
                    </div>
                </div>
            </div>
        </div>
    )
}