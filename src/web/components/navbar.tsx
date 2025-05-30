"use client"

import clsx from "clsx";
import {LinkTo, PageName} from "@/components/pageEnums";
import {togglePagination} from "@/states/paginationStateHandler";
import {redirect, usePathname} from "next/navigation";


export default function Navbar() {
    const resetPagination = togglePagination((state) => state.reset)
    const currentPageName = usePathname() + ".html"

    const links = [
        {name: PageName.MANAGEMENT, href: LinkTo.MANAGEMENT},
        {name: PageName.MAILINGLISTS, href: LinkTo.MAILINGLISTS},
        {name: PageName.APPROVALREQUESTS, href: LinkTo.APPROVALREQUESTS},
    ];

    const handleLogout = () => {
        localStorage.removeItem("sessionToken");
        redirect(LinkTo.LOGIN);
    }

    return (
        <div className="
        flex
        h-full
        flex-col
        py-4
        shadow-[0_3px_10px_-1px_rgba(0,0,0,1)]
        ">
            {links.map((link) => (
                    <a href={`/ui${link.href}`} key={link.name} data-testid={link.name} onClick={() => resetPagination()}
                          className={clsx(
                              "flex h-[60px] border-t-1 border-b-1 border-neutral-700  shadow-[0_3px_10px_-1px_rgba(0,0,0,0.3)] grow items-center justify-center gap-2 p-3 text-sm font-medium hover:bg-cyan-500 hover:text-gray-200 hover:cursor-pointer md:flex-none md:justify-start w-full",
                              {
                                  "bg-cyan-600 border-t-1 border-b-1 border-neutral-700 shadow-[0_3px_10px_-1px_rgba(0,0,0,0.3)]  text-gray-200 font-bold": currentPageName === link.href
                              }
                          )}>{link.name}</a>
                )
            )}

            <button
                onClick={handleLogout}
                data-testid={PageName.LOGIN}
                className="flex h-[60px] shadow-[0_3px_10px_-1px_rgba(0,0,0,0.3)] grow items-center justify-center gap-2 rounded-md p-3 text-sm font-medium hover:bg-red-700 hover:text-gray-200 hover:cursor-pointer md:flex-none md:justify-start w-full">
                Sign Out
            </button>
        </div>
    )
}
