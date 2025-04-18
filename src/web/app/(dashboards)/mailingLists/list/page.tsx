"use client"

import Pagination from "@/components/pagination";
import {Suspense} from "react";
import ListEmailsTableWrapper from "@/components/list/listEmailsTableWrapper";

export default function Page() {


    return (
        <div className="w-full">
            <Suspense>
                <ListEmailsTableWrapper />
            </Suspense>

            <Suspense>
                <Pagination totalPages={1} />
            </Suspense>


        </div>

    )
}