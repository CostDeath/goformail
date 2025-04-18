"use client"

import Pagination from "@/components/pagination";
import {Suspense} from "react";
import ListEmailsTable from "@/components/list/listEmailsTable";

export default function Page() {


    return (
        <div className="w-full">
            <Suspense>
                <ListEmailsTable />
            </Suspense>

            <Suspense>
                <Pagination totalPages={1} />
            </Suspense>


        </div>

    )
}