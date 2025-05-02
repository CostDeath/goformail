"use client"

import EmailApprovalForm from "@/components/emailApprovalRequests/emailApprovalForm";
import {Suspense} from "react";

export default function Page() {
    return (
        <>
            <Suspense>
                <EmailApprovalForm />
            </Suspense>
        </>
    )
}