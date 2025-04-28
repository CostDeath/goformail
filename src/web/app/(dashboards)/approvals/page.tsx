"use client"

import ApprovalRequestsTable from "@/components/emailApprovalRequests/approvalRequestsTable";
import Pagination from "@/components/pagination";
import {Suspense} from "react";

export default function Page() {
    return (
        <div className="w-full">
            <ApprovalRequestsTable api="placeholder" />
            <Suspense>
                <Pagination totalPages={1} />
            </Suspense>


        </div>
    )
}