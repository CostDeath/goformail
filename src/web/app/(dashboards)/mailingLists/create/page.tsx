"use client"

import {ChangeEvent, useState} from "react";
import StyledInput from "@/components/styledInput";


export default function Page() {
    const [name, setName] = useState("");
    const [recipients, setRecipients] = useState([{value: ""}])

    const placeholder = () => {
        console.log(name)
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
        if (recipients.length > 1) {
            const values = [...recipients]
            values.splice(index, 1)
            setRecipients(values)
        }
    }

    return (
        <>
            <h1 className="font-bold px-2 text-2xl">Create a new mailing list</h1>
            <div className="grid grid-cols-2 py-10">
                <div className="px-5 text-xl">Mailing List Name</div>
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
                    type="text"
                    name="listName"
                    value={name}
                    onChange={e => setName(e.target.value)}
                    required
                />
            </div>

            <div className="flex flex-row justify-end px-5">
                <button className="bg-green-600/75 hover:bg-green-600 px-2 py-1 rounded-md" onClick={placeholder}>Submit</button>
            </div>
        </>
    )
}