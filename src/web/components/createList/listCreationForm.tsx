"use client"

import {ChangeEvent, useState} from "react";
import {api} from "@/components/api";
import {redirect} from "next/navigation";
import {LinkTo} from "@/components/pageEnums";
import validateEmail from "@/components/validateEmails";

export default function ListCreationForm() {
    const [name, setName] = useState("");
    const [recipients, setRecipients] = useState([{value: ""}])

    const createList = async () => {
        const url = `${window.location.origin}/api${api.list}`
        if (!validateEmail(name)) {
            alert("Please enter a valid mailing list email")
            return;
        }
        const rcptList: string[] = []
        for (let i = 0; i < recipients.length; i++) {
            if (!validateEmail(recipients[i].value)) {
                alert(`${recipients[i].value} is not a valid email`)
                return;
            }
            rcptList.push(recipients[i].value)
        }

        const response = await fetch(url, {
            method: "POST",
            body: JSON.stringify({
                Name: name,
                Recipients: rcptList
            }),
            headers: {
                "Content-Type": "application/json"
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

    return (
        <>
            <div className="grid grid-cols-2 py-10">
                <label htmlFor="listName" className="px-5 text-xl">Mailing List Email</label>
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
                    placeholder="Mailing List Email"
                    onChange={e => setName(e.target.value)}
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
                        onClick={createList}>Submit
                </button>
            </div>
        </>
    )
}