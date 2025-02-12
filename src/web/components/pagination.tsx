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
                        "flex h-10 w-10 items-center justify-center text-sm border rounded-l-md",
                        {
                            "z-10 bg-blue-600 border-blue-600 text-white": currentPage === key,
                            "hover: bg-gray-300": currentPage !== key
                        }
                    )}>{key}
                    </div>
                </div>
            ))}
        </div>
    )

}