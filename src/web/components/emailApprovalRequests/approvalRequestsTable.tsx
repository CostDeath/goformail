//import useSWR from "swr";


export default function ApprovalRequestsTable({currentPage, api}: {
    currentPage?: number;
    api: string;
}) {
    // let data = await fetch(api);
    // let retriedData = data.json();
    console.log(api)
    console.log(currentPage)
    // const [id, setId] = useState("");

    /* potentially what we can use to fetch
    const fetcher = (...args) => fetch(...args).then((res) => res.json())

    const {data, error} = useSWR("https://catfact.ninja/fact", fetcher)

    if (error) return <div>Failed to load</div>
    if (!data) return <div>Loading...</div>

     */


    return (
        <>
        <div className="min-w-full table text-gray-900 shadow-[0_3px_10px_-1px_rgba(0,0,0,1)]">
            <div className="table-header-group text-left text-sm font-normal">
                <div data-testid="table-head" className="table-row bg-neutral-800/45 text-neutral-300">
                    <div className="table-cell border-b border-black text-left px-4 py-5 font-bold sm:pl-6">
                        Mailing List Email {/*data.fact*/}
                    </div>
                    <div className="table-cell border-b border-black text-left px-4 py-5 font-bold sm:pl-6">
                        Title
                    </div>
                    <div className="table-cell border-b border-black text-left px-4 py-5 font-bold sm:pl-6">
                        Date Sent for Approval
                    </div>
                </div>
            </div>
            <div data-testid="table-body" className="table-row-group">
                <a href={`/ui/approvals/email.html?id=${1}`} className="table-row shadow-inner text-neutral-300 hover:bg-neutral-600/75  hover:cursor-pointer">
                    <div className="table-cell border-black border-b py-3 text-sm">
                        <div className="whitespace-nowrap py-3 pl-6 pr-3 flex items-center gap-3">
                            exampleentry@email.com
                        </div>
                    </div>
                    <div className="table-cell border-black border-b py-3 text-sm">
                        <div className="whitespace-nowrap py-3 pl-6 pr-3 flex items-center gap-3">
                            Testing Forwarding
                        </div>
                    </div>
                    <div className="table-cell border-black border-b py-3 text-sm">
                        <div className="whitespace-nowrap py-3 pl-6 pr-3 flex items-center gap-3">
                            09/01/2025
                        </div>
                    </div>
                </a>

                <a href={`/ui/approvals/email.html?id=${2}`} className="table-row shadow-inner text-neutral-300 hover:bg-neutral-600/75  hover:cursor-pointer">
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