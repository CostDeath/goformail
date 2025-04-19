"use client"

import ListCreationForm from "@/components/createList/listCreationForm";


export default function Page() {
    return (
        <>
            <h1 className="font-bold px-2 text-2xl">Create a new mailing list</h1>
            <ListCreationForm />
        </>
    )
}