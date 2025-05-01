"use client"

import {redirect, useSearchParams} from "next/navigation";
import {ChangeEvent, useEffect, useState} from "react";
import {getSessionToken} from "@/components/sessionToken";
import useSWR from "swr";
import {api} from "@/components/api";
import {MailingList} from "@/models/list";
import validateEmail from "@/components/validateEmails";

export default function ManageRecipientsForm() {
    const params = useSearchParams()
    const listId = params.get("id")
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
    } else if (data.message !== "Successfully fetched list!") return <div>Error</div>

    const result: MailingList = data.data

    const editRecipients = async () => {
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
                recipients: rcptList,
            }),
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${sessionToken}`
            }
        })

        if (response.ok) {
            alert("Successfully edited recipients!")
            redirect(`/mailingLists/list.html?id=${listId}`)
        } else {
            const result = await response.text()
            alert(result)
        }
    }

    return (
        <>
            <h1 className="text-2xl font-bold px-2 py-5">{result.name}</h1>
            <hr/>
            <h1 className="px-2 py-5 text-2xl underline">Manage Recipients</h1>
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
                            required/>
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
                        onClick={handleAdd}>
                    + Add Another Recipient
                </button>
            </div>

            <div className="flex flex-row justify-end px-5">
                <button className="bg-green-600/75 hover:bg-green-600 px-2 py-1 rounded-md"
                        onClick={editRecipients}>Submit
                </button>
            </div>
        </>
    )
}