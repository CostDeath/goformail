import {redirect, useSearchParams} from "next/navigation";
import {ChangeEvent, useEffect, useState} from "react";
import useSWR from "swr";
import DeleteList from "@/components/editList/deleteList";
import {api} from "@/components/api";
import {List} from "@/models/list";
import validateEmail from "@/components/validateEmails";
import {LinkTo} from "@/components/pageEnums";
import {getSessionToken} from "@/components/sessionToken";

export default function ListEditForm() {
    const searchParams = useSearchParams()
    const listId = searchParams.get("id")
    const [locked, setLocked] = useState(false)
    const [recipients, setRecipients] = useState([{value: ""}])
    const [sessionToken, setSessionToken] = useState<string | null>()


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
        const response = await fetch(url, {
            method: "GET",
            headers: {
                "Authorization": `Bearer ${sessionToken}`
            }
        })
        const data = await response.json()
        const rcptsBuilder: { value: string }[] = []
        if (data.data.recipients) {
            const rcpts: string[] = data.data.recipients
            for (let i = 0; i < rcpts.length; i++) {
                rcptsBuilder.push({value: rcpts[i]})
            }
        }

        setRecipients(rcptsBuilder)
        setLocked(data.data.locked)
        return data
    }

    const [baseUrl, setBaseUrl] = useState("")

    useEffect(() => {
        const url =`${window.location.origin}/api`
        setBaseUrl(url)
        setSessionToken(getSessionToken())
    }, [])

    const {data, error} = useSWR((baseUrl && sessionToken) ? `${baseUrl}${api.list}?id=${listId}` : null, fetcher)

    if (error) return <div>Error</div>
    if (!data) {
        return <div>Loading</div>
    } else if (data.message !== "Successfully fetched list!") {
        return <div>Error</div>
    }

    const result: List = data.data


    const editList = async () => {
        const url = `${window.location.origin}/api${api.list}?id=${listId}`
        const rcptList: string[]  = []
        for (let i = 0; i < recipients.length; i++) {
            if (!validateEmail(recipients[i].value)) {
                alert(`${recipients[i].value} is not a valid email`)
                return;
            }
            rcptList.push(recipients[i].value)
        }

        const response = await fetch(url, {
            method: "PATCH",
            body: JSON.stringify({
                Name: result.name,
                Recipients: rcptList,
                Locked: locked
            }),
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${sessionToken}`
            }
        })

        if (response.ok) {
            const result = await response.json()
            alert(result.message)
            redirect(LinkTo.MAILINGLISTS)
        } else {
            const result = await response.text()
            alert(result)
        }
    }

    return (
        <>
            <DeleteList id={listId}/>
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
                    value={result.name}
                    disabled
                />
            </div>

            <div className="grid grid-cols-2 py-5">
                <label htmlFor="locked" className="px-5 text-xl">Mailing List Locked?</label>
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
                    checked={locked}
                    id="locked"
                    type="checkbox"
                    name="locked"
                    value=""
                    onChange={e => setLocked(e.target.checked)}
                    required
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
                            <button aria-label={`delete${index}`}
                                    className="bg-red-600 hover:bg-red-700 py-2 px-3 rounded-md"
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
                        onClick={editList}>Submit
                </button>
            </div>
        </>
    )
}