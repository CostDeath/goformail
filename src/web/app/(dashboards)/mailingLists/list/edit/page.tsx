"use client"

import ListEditForm from "@/components/editList/listEditForm";
import {Suspense} from "react";
import Card from "@/components/card";


export default function Page() {
    return (
        <>
            <Card>
            <Suspense>
                <ListEditForm />
            </Suspense>
            </Card>
        </>
    )
}
