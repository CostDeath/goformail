"use client"

import ApprovalRequestsTable from "@/components/emailApprovalRequests/approvalRequestsTable";
import Pagination from "@/components/pagination";
import Modal from "@/components/modal";
import EmailView from "@/components/emailApprovalRequests/emailView";
import {useModal} from "@/states/modalStateHandler";
import {Suspense} from "react";

export default function ApprovalRequestsPage() {
    const showModal = useModal((state) => state.toggled)
    return (
        <div className="w-full">
            <ApprovalRequestsTable api="placeholder" />
            <Suspense>
            <Pagination totalPages={1} />
            </Suspense>

            {showModal && (
                <Modal><EmailView /></Modal>
            )}
        </div>
    )
}