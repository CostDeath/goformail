"use client"

import {useSearchParams} from "next/navigation";
import useSWR from "swr";
import {useModal} from "@/states/modalStateHandler";
import ListEmailsTable from "@/components/list/listEmailsTable";
import Modal from "@/components/modal";
import EmailView from "@/components/emailApprovalRequests/emailView";
import Pagination from "@/components/pagination";
import {Suspense} from "react";

export default function Page() {
    const showModal = useModal((state) => state.toggled)
    const searchParams = useSearchParams()
    const id = searchParams.get("id")

    const fetcher = (...args: Parameters<typeof fetch>) =>
        fetch(...args).then((res) => res.json())

    const {data, error} = useSWR(`https://jsonplaceholder.typicode.com/posts/${id}`, fetcher)

    if (error) return <div>Error</div>
    if (!data) {
        return <div>Loading...</div>
    } else if (!data.id) {
        return <div>Error</div>
    }


    console.log(data)

    return (
        <div className="w-full">
            <h1 className="font-bold py-5 px-2 text-2xl">{data.title}</h1>
            <ListEmailsTable/>

            <Suspense>
                <Pagination totalPages={1} />
            </Suspense>

            {showModal && (<Modal><EmailView /></Modal>)}

        </div>

    )
}