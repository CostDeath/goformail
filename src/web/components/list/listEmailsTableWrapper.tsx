import ListEmailsTable from "@/components/list/listEmailsTable";
import {useSearchParams} from "next/navigation";
import useSWR from "swr";


export default function ListEmailsTableWrapper() {
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
        <>
            <h1 className="font-bold py-5 px-2 text-2xl">{data.title}</h1>
            <ListEmailsTable/>
        </>
    )
}