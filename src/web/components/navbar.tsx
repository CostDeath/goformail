"use client"

import clsx from "clsx";
import {PageName, switchPage} from "@/states/linkStateHandler";
import MailingListsPage from "@/app/pages/mailingListsPage";
import ApprovalRequestsPage from "@/app/pages/approvalRequestsPage";
import LoginSignupPage from "@/app/pages/LoginSignupPage";
import {togglePagination} from "@/states/paginationStateHandler";


export default function Navbar() {
    const currentPageName = switchPage((state) => state.name);
    const changePage = switchPage((state) => state.changePage);
    const resetPagination = togglePagination((state) => state.reset)

    const links = [
        {name: PageName.MAILINGLISTS, href: () => {changePage(<MailingListsPage />, PageName.MAILINGLISTS); resetPagination()}},
        {name: PageName.APPROVALREQUESTS, href: () => {changePage(<ApprovalRequestsPage />, PageName.APPROVALREQUESTS); resetPagination()}},
    ];

    return (
        <div className="flex h-full flex-col px-3 py-4 md:px-2 bg-gray-100">
            {links.map((link) => (
                <div onClick={link.href} key={link.name} data-testid={link.name}
                className={clsx(
                    "flex h-[48px] border-b-4 grow items-center justify-center gap-2 rounded-md p-3 text-sm font-medium hover:bg-gray-200 hover:text-gray-700 hover:cursor-pointer md:flex-none md:justify-start w-full",
                    {
                        "bg-gray-200 text-gray-700 font-bold": currentPageName === link.name
                    }
                )}>{link.name}</div>
            ))}

            <div data-testid={PageName.LOGINSIGNUP} onClick={() => {changePage(<LoginSignupPage/>, PageName.LOGINSIGNUP); resetPagination()}} className="flex h-[48px] border-b-4 grow items-center justify-center gap-2 rounded-md p-3 text-sm font-medium hover:bg-gray-200 hover:text-gray-700 hover:cursor-pointer md:flex-none md:justify-start w-full">
                {/* For now we'll use Link to sign out as there is no logic for users yet */}
                Sign Out
            </div>
        </div>
    )
}
