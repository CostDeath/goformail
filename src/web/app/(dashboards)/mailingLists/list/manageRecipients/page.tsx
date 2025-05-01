"use client"

import ManageRecipientsForm from "@/components/manageRecipients/manageRecipientsForm";
import Card from "@/components/card";
import {Suspense} from "react";

export default function Page() {
    return (
        <>
            <Card>
                <Suspense>
                    <ManageRecipientsForm />
                </Suspense>
            </Card>
        </>
    )
}