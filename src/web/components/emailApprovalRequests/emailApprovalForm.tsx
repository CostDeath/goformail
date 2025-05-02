"use client"

import {redirect, useSearchParams} from "next/navigation";
import {api} from "@/components/api";
import {getSessionToken} from "@/components/sessionToken";
import {LinkTo} from "@/components/pageEnums";
import {useEffect, useState} from "react";
import useSWR from "swr";
import {Email} from "@/models/email";


export default function EmailApprovalForm() {
    const search = useSearchParams()
    const id = search.get("id")

    const approveEmail = async () => {
        const url = `${window.location.origin}/api${api.emails}approve/?id=${id}`
        const sessionToken = getSessionToken()

        const response = await fetch(url, {
            method: "POST",
            body: JSON.stringify({}),
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

    const [sessionToken, setSessionToken] = useState<string | null>()
    const [baseUrl, setBaseUrl] = useState("")

    useEffect(() => {
        const url = `${window.location.origin}/api`
        setBaseUrl(url)
        setSessionToken(getSessionToken())
    }, [])

    const fetcher = async (url: string) => {
        const response = await fetch(url, {
            method: "GET",
            headers: {
                "Authorization": `Bearer ${sessionToken}`
            }
        })
        if (response.ok) {
            return await response.json()
        }
        return response.text()
    }

    const {data, error} = useSWR((baseUrl && sessionToken) ? `${baseUrl}${api.email}?id=${id}` : null, fetcher)

    if (error) {
        return <div>Error</div>
    }
    if (!data) {
        return <div>Loading</div>
    } else if (data.message !== "Successfully fetched email!") {
        return <div>Error</div>
    }

    const email: Email = data.data
    const rcpts = email.rcpt.join(", ")

    return (
        <>

            <div className="">
                <div data-testid="email-title" className="border-b border-black py-2">To: {rcpts}</div>
                <div data-testid="email-subject" className="border-b border-black py-2">From: {email.sender}</div>
            </div>
            <div data-testid="email-content" className="border border-black py-5 px-4 min-h-[50vh] bg-neutral-700">
                {email.content}
            </div>


            <div className="flex flex-row justify-end py-6 px-3">
                <div className="p-2">
                    <button className="bg-green-600 text-white p-2 rounded-xl hover:bg-green-700"
                            onClick={approveEmail}>Approve
                    </button>
                </div>
            </div>
        </>
    )
}