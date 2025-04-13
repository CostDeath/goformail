"use client"

import clsx from "clsx";
import {togglePagination} from "@/states/paginationStateHandler";

export default function Pagination({totalPages}: { totalPages: number }) {
    const currentPage = togglePagination((state) => state.pageNumber);
    const pageToggler = togglePagination((state) => state.changePage);
    const allPages = [];

    for (let i = 1; i < totalPages + 1; i++) {
        allPages.push(i)
    }

    return (
        <div data-testid="pagination" className="flex justify-center py-5">
            {allPages.map(key => (
                <div key={key} onClick={() => pageToggler(key)} className="hover:cursor-pointer">
                    <div className={clsx(
                        "flex h-10 w-10 items-center justify-center text-sm border border-neutral-700/25 rounded-md",
                        {
                            "z-10 bg-cyan-600 border-cyan-600 text-white": currentPage === key,
                            "hover: bg-neutral-700 hover:bg-neutral-600": currentPage !== key
                        }
                    )}>{key}
                    </div>
                </div>
            ))}
        </div>
    )

}