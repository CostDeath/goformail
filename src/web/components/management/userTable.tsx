import Link from "next/link";

export default function UserTable() {
    return(
        <>
            <div className="min-w-full table text-gray-900 shadow-[0_3px_10px_-1px_rgba(0,0,0,1)]">
                <div className="table-header-group text-left text-sm font-normal">
                    <div data-testid="table-head" className="table-row bg-neutral-800/45 text-neutral-300">
                        <div className="table-cell border-b border-black text-left px-4 py-5 font-bold sm:pl-6">
                            User
                        </div>
                    </div>
                </div>
                <div data-testid="table-body" className="table-row-group">
                    <Link href={`/management/edit?id=${1}`} className="table-row shadow-inner text-neutral-300 hover:bg-neutral-600/75  hover:cursor-pointer">
                        <div className="table-cell border-black border-b py-3 text-sm">
                            <div className="whitespace-nowrap py-3 pl-6 pr-3 flex items-center gap-3">
                                user@example.com
                            </div>
                        </div>
                    </Link>

                    <Link href={`/management/edit?id=${2}`} className="table-row shadow-inner text-neutral-300 hover:bg-neutral-600/75  hover:cursor-pointer">
                        <div className="table-cell black border-black py-3 text-sm">
                            <div className="whitespace-nowrap py-3 pl-6 pr-3 flex items-center gap-3">
                                user2@example.com
                            </div>
                        </div>
                    </Link>
                </div>
            </div>
        </>
    )
}