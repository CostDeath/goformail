"use client"

import ManageModsForm from "@/components/manageMods/manageModsForm";
import Card from "@/components/card";
import {Suspense} from "react";

export default function Page() {
    return (
        <>
            <Card>
                <Suspense>
                    <ManageModsForm />
                </Suspense>
            </Card>
        </>
    )
}