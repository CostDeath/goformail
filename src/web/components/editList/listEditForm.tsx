import {useSearchParams} from "next/navigation";
import {ChangeEvent, useState} from "react";
import useSWR from "swr";
import DeleteList from "@/components/editList/deleteList";

export default function ListEditForm() {
    const searchParams = useSearchParams()
    const listId = searchParams.get("id")
    const [recipients, setRecipients] = useState([{value: ""}])


    const handleChange = (index: number, e: ChangeEvent<HTMLInputElement>) => {
        const values = [...recipients]
        values[index].value = e.target.value
        setRecipients(values)
    }

    const handleAdd = () => {
        setRecipients([...recipients, {value: ""}])
    }

    const handleRemove = (index: number) => {
        const values = [...recipients]
        values.splice(index, 1)
        setRecipients(values)
    }


    const fetcher = async(url: string) => {
        const response = await fetch(url)
        const data = await response.json()
        // placeholder
        setRecipients([{value: data.title}])
        return data
    }

    const {data, error} = useSWR(`https://jsonplaceholder.typicode.com/posts/${listId}`, fetcher)

    if (error) return <div>Error</div>
    if (!data) {
        return <div>Loading</div>
    } else if (!data.id) {
        return <div>Error</div>
    }

    // TODO: uncomment this once doing ticket for connecting frontend and backend
    // setRecipients(data.Recipients) // Assumption that this is how it'll work when fetched


    const placeholder = () => {
        // This will be a patch request
        console.log(data)
        console.log(recipients)
    }

    return (
        <>
            <DeleteList />
            <div className="grid grid-cols-2 py-10">
                <label htmlFor="listName" className="px-5 text-xl">Mailing List Name</label>
                <input
                    className="
                bg-neutral-700
                peer
                block
                w-full
                h-10
                px-3
                border
                border-neutral-500
                rounded-md
                outline-2
                placeholder:text-neutral-500
                "
                    id="listName"
                    type="email"
                    name="listName"
                    value={data.Name}
                    disabled
                />
            </div>
            <br/>
            <hr/>
            <br/>
            <h1 className="px-2 text-2xl underline">Add recipients</h1>
            <div className="py-10">
                {recipients.map((recipient, index) => (
                    <div className="grid grid-cols-3 px-2 py-4" key={index}>
                        <input
                            className="
                bg-neutral-700
                peer
                block
                w-full
                h-10
                px-3
                border
                border-neutral-500
                rounded-md
                outline-2
                placeholder:text-neutral-500
                "
                            id={`recipient${index}`}
                            type="email"
                            name={`recipient${index}`}
                            aria-label={`recipient${index}`}
                            value={recipient.value}
                            onChange={e => handleChange(index, e)}
                            placeholder="Email Address"
                            required
                        />
                        <div className="px-7">
                            <button aria-label={`delete${index}`} className="bg-red-600 hover:bg-red-700 py-2 px-3 rounded-md"
                                    onClick={() => handleRemove(index)}>
                                Remove
                            </button>
                        </div>
                    </div>

                ))}
                <button className="bg-cyan-600 text-white hover:bg-cyan-500 py-2 px-3 rounded-md font-bold"
                        onClick={handleAdd}>+ Add recipient
                </button>
            </div>

            <div className="flex flex-row justify-end px-5">
                <button className="bg-green-600/75 hover:bg-green-600 px-2 py-1 rounded-md"
                        onClick={placeholder}>Submit
                </button>
            </div>
        </>
    )
}