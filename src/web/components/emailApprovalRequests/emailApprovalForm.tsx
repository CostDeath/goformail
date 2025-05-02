"use client"

import EmailView from "@/components/emailView";
import {redirect, useSearchParams} from "next/navigation";
import {api} from "@/components/api";
import {getSessionToken} from "@/components/sessionToken";
import {LinkTo} from "@/components/pageEnums";
// import {useState} from "react";


export default function EmailApprovalForm() {
    /*
    //TODO: Remove this post example once doing fetching ticket
    const [title, setTitle] = useState("");
    const [body, setBody] = useState("");

    const submitApproval = async () => {
        let response = await fetch('https://jsonplaceholder.typicode.com/posts', {
            method: 'POST',
            body: JSON.stringify({
                title: title,
                body: body,
                userId: 1
            }),
            headers: {
                "Content-Type": "application/json"
            }
        })

        response = await response.json();
        alert(JSON.stringify(response))
    }

     */

    const search = useSearchParams()
    const id = search.get("id")

    const approveEmail = async () => {
        const url = `${window.location.origin}${api.emails}${api.approveEmail}?id=${id}`
        const sessionToken = getSessionToken()

        const response = await fetch(url, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${sessionToken}`,
            }
        })
        if (response.ok) {
            const data = await response.json()
            alert(data.message)
            redirect(LinkTo.APPROVALREQUESTS)
        } else {
            const message = await response.text()
            alert(message)
        }

    }

    if (!id || isNaN(Number(id))) {
        return <div>Error</div>
    }

    // TODO: Figure out how to send data to email view while also hiding approve button if no email was found

    return (
        <>
            <EmailView />
            <div className="flex flex-row justify-end py-6 px-3">
                <div className="p-2">
                    <button className="bg-green-600 text-white p-2 rounded-xl hover:bg-green-700" onClick={approveEmail}>Approve</button>
                </div>
            </div>


            {/*
            <br/>
            <br/>
            <h1>This is a post test</h1>
            <div>
                <h2>Testing posts</h2>
                <input type="text" value={title} onChange={e => setTitle(e.target.value)}/>
                <input type="text" value={body} onChange={e => setBody(e.target.value)}/>

                <button onClick={submitApproval}>Submit</button>
            </div>
            */}
        </>
    )
}