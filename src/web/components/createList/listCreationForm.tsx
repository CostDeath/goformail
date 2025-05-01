"use client"

import {ChangeEvent, useState} from "react";
import {api} from "@/components/api";
import {redirect} from "next/navigation";
import validateEmail from "@/components/validateEmails";
import {getSessionToken} from "@/components/sessionToken";

export default function ListCreationForm() {
    const [name, setName] = useState("");
    const [senders, setSenders] = useState([{value: ""}])
    const [locked, setLocked] = useState(false)

    const createList = async () => {
        const url = `${window.location.origin}/api${api.list}`
        const token = getSessionToken()
        const senderList: string[] = []
        for (let i = 0; i < senders.length; i++) {
            if (!senders[i].value) {
                alert("Cannot leave sender input blank, either remove the sender or fill in a valid email")
            } else if (!validateEmail(senders[i].value)) {
                alert(`${senders[i].value} is not a valid email`)
                return;
            }
            senderList.push(senders[i].value)
        }

        const response = await fetch(url, {
            method: "POST",
            body: JSON.stringify({
                Name: name,
                Recipients: [],
                locked: locked,
                Mods: [],
                approved_senders: senderList
            }),
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${token}`
            }
        })

        if (response.ok) {
            const result = await response.json()
            alert(result.message)
            redirect(`/mailingLists/list/?id=${result.data.id}`)
        } else {
            const result = await response.text()
            alert(result)
        }
    }

    const handleChange = (index: number, e: ChangeEvent<HTMLInputElement>) => {
        const values = [...senders]
        values[index].value = e.target.value
        setSenders(values)
    }

    const handleAdd = () => {
        setSenders([...senders, {value: ""}])
    }

    const handleRemove = (index: number) => {
        const values = [...senders]
        values.splice(index, 1)
        setSenders(values)
    }

    return (
        <>
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
                    value={name}
                    placeholder="Mailing List Email Name"
                    onChange={e => setName(e.target.value)}
                    required
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
            <h1 className="px-2 text-2xl underline">Add Approved Senders</h1>
            <div className="py-10">
                {senders.map((sender, index) => (
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
                            id={`sender${index}`}
                            type="email"
                            name={`sender${index}`}
                            aria-label={`sender${index}`}
                            value={sender.value}
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
                        onClick={handleAdd}>+ Add Another Sender
                </button>
            </div>

            <div className="flex flex-row justify-end px-5">
                <button className="bg-green-600/75 hover:bg-green-600 px-2 py-1 rounded-md"
                        onClick={createList}>Submit
                </button>
            </div>
        </>
    )
}