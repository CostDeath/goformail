"use client"

import {useSearchParams} from "next/navigation";
import useSWR from "swr";
import ListEmailsTable from "@/components/list/listEmailsTable";
import Pagination from "@/components/pagination";
import {Suspense} from "react";

export default function Page() {
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


        </div>

    )
}