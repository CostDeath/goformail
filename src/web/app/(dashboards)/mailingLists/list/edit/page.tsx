"use client"

import ListEditForm from "@/components/editList/listEditForm";
import {Suspense} from "react";


export default function Page() {
    return (
        <>
            <Suspense>
                <ListEditForm />
            </Suspense>
        </>
    )
}
