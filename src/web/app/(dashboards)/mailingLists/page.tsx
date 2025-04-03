import MailingListTable from "@/components/mailingLists/malingListTable";
import Pagination from "@/components/pagination";
import {Suspense} from "react";

export default function Page() {
    return (
        <div className="w-full">
            <MailingListTable api="placeholder" />
            <Suspense>
                <Pagination totalPages={2} />
            </Suspense>
        </div>
    )
}