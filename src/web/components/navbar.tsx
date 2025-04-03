"use client"

import clsx from "clsx";
import {LinkTo, PageName} from "@/components/pageEnums";
import {togglePagination} from "@/states/paginationStateHandler";
import Link from "next/link";
import {usePathname} from "next/navigation";


export default function Navbar() {
    const resetPagination = togglePagination((state) => state.reset)
    const currentPageName = usePathname()

    const links = [
        {name: PageName.MAILINGLISTS, href: "mailingLists"},
        {name: PageName.APPROVALREQUESTS, href: "approvals"},
    ];

    return (
        <div className="flex h-full flex-col px-3 py-4 md:px-2 bg-gray-100">
            {links.map((link) => (
                <Link href={link.href} key={link.name} data-testid={link.name} onClick={() => resetPagination()}
                className={clsx(
                    "flex h-[48px] border-b-4 grow items-center justify-center gap-2 rounded-md p-3 text-sm font-medium hover:bg-gray-200 hover:text-gray-700 hover:cursor-pointer md:flex-none md:justify-start w-full",
                    {
                        "bg-gray-200 text-gray-700 font-bold": currentPageName === link.href
                    }
                )}>{link.name}</Link>
            ))}

            <Link
                href={LinkTo.LOGIN}
                data-testid={PageName.LOGIN}
                 className="flex h-[48px] border-b-4 grow items-center justify-center gap-2 rounded-md p-3 text-sm font-medium hover:bg-gray-200 hover:text-gray-700 hover:cursor-pointer md:flex-none md:justify-start w-full">
                {/* For now we'll use Link to sign out as there is no logic for users yet */}
                Sign Out
            </Link>
        </div>
    )
}
