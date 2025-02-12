import {useModal} from "@/states/modalStateHandler";
//import useSWR from "swr";


export default function ApprovalRequestsTable({currentPage, api}: {
    currentPage?: number;
    api: string;
}) {
    // having api be nullable for now, TODO: make it not nullable once api is ready
    // let data = await fetch(api);
    // let retriedData = data.json();
    console.log(api)
    console.log(currentPage)

    const toggleModal = useModal((state) => state.toggleModal);

    /* potentially what we can use to fetch
    const fetcher = (...args) => fetch(...args).then((res) => res.json())

    const {data, error} = useSWR("https://catfact.ninja/fact", fetcher)

    if (error) return <div>Failed to load</div>
    if (!data) return <div>Loading...</div>

     */


    return (
        <div className="min-w-full table text-gray-900 border-gray-400 border-[1px]">
            <div className="table-header-group text-left text-sm font-normal">
                <div data-testid="table-head" className="table-row">
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
                <div className="table-row shadow-inner hover:bg-gray-100 hover:cursor-pointer" onClick={() => toggleModal(true)}>
                    <div className="table-cell border-gray-200 border-b py-3 text-sm">
                        <div className="whitespace-nowrap py-3 pl-6 pr-3 flex items-center gap-3">
                            exampleentry@email.com
                        </div>
                    </div>
                    <div className="table-cell border-gray-200 border-b py-3 text-sm">
                        <div className="whitespace-nowrap py-3 pl-6 pr-3 flex items-center gap-3">
                            Testing Forwarding
                        </div>
                    </div>
                    <div className="table-cell border-gray-200 border-b py-3 text-sm">
                        <div className="whitespace-nowrap py-3 pl-6 pr-3 flex items-center gap-3">
                            09/01/2025
                        </div>
                    </div>
                </div>

                <div className="table-row shadow-inner hover:bg-gray-100 hover:cursor-pointer" onClick={() => toggleModal(true)}>
                    <div className="table-cell border-gray-200 border-b py-3 text-sm">
                        <div className="whitespace-nowrap py-3 pl-6 pr-3 flex items-center gap-3">
                            exampleentry2@email.com
                        </div>
                    </div>
                    <div className="table-cell border-gray-200 border-b py-3 text-sm">
                        <div className="whitespace-nowrap py-3 pl-6 pr-3 flex items-center gap-3">
                            Testing Testing...
                        </div>
                    </div>
                    <div className="table-cell border-gray-200 border-b py-3 text-sm">
                        <div className="whitespace-nowrap py-3 pl-6 pr-3 flex items-center gap-3">
                            09/01/2025
                        </div>
                    </div>
                </div>
            </div>
        </div>
    )
}