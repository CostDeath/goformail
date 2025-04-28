import MailingListTable from "@/components/mailingLists/malingListTable";
import Pagination from "@/components/pagination";
import {Suspense} from "react";
import Link from "next/link";

export default function Page() {
    return (
        <div className="w-full">
            <div className="py-5">
            <Link href="/mailingLists/create.html" className="bg-cyan-600 text-white hover:bg-cyan-500 px-3 py-2 rounded-md font-bold">+ New Mailing List</Link>
            </div>
            <MailingListTable api="placeholder" />
            <Suspense>
                <Pagination totalPages={2} />
            </Suspense>
        </div>
    )
}