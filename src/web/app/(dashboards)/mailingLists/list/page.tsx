"use client"

import {Suspense} from "react";
import ListEmailsTable from "@/components/list/listEmailsTable";
import RecipientsTable from "@/components/list/recipientsTable";

export default function Page() {


    return (
        <div className="w-full">
            <Suspense>
                <ListEmailsTable />
            </Suspense>

            <br />

            <Suspense>
                <RecipientsTable />
            </Suspense>
        </div>

    )
}