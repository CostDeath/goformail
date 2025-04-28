"use client"

import ListCreationForm from "@/components/createList/listCreationForm";
import Card from "@/components/card";


export default function Page() {
    return (
        <>
            <Card>
            <h1 className="font-bold px-2 py-5 text-2xl">Create a new mailing list</h1>
                <hr/>
            <ListCreationForm />
            </Card>
        </>
    )
}