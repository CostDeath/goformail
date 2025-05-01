"use client"

import {useSearchParams} from "next/navigation";
import {useEffect, useState} from "react";
import {getSessionToken} from "@/components/sessionToken";
import useSWR from "swr";
import {api} from "@/components/api";
import {User} from "@/models/user";
import ModeratorTable from "@/components/manageMods/moderatorTable";
import ModUserTable from "@/components/manageMods/modUserTable";

export default function ManageModsForm() {
    const params = useSearchParams()
    const listId = params.get("id")
    const [baseUrl, setBaseUrl] = useState("")
    const [listName, setListName] = useState("")
    const [moderators, setModerators] = useState<{id: number, email: string}[]>([])
    const [modList, setModList] = useState<number[]>([])
    const [sessionToken, setSessionToken] = useState<string | null>()

    const fetcher = async(url: string) => {
        let response = await fetch(url, {
            method: "GET",
            headers: {
                "Authorization": `Bearer ${sessionToken}`
            }
        })
        const listData = await response.json()
        setListName(listData.data.name)
        const mods = listData.data.mods
        setModList(mods)

        const modBuilder: {id: number; email: string}[] = []
        let user: User
        if (mods) {
            for (let i = 0; i < mods.length; i++) {
                response = await fetch(`${baseUrl}${api.user}?id=${mods[i]}`, {
                    method: "GET",
                    headers: {
                        "Authorization": `Bearer ${sessionToken}`
                    }
                })
                if (response.ok) {
                    const userData = await response.json()
                    user = userData.data
                    modBuilder.push({id: user.id, email: user.email})
                }
            }
        }
        setModerators(modBuilder)

        response = await fetch(`${baseUrl}${api.users}`, {
            method: "GET",
            headers: {
                "Authorization": `Bearer ${sessionToken}`
            }
        })
        return await response.json()
    }

    useEffect(() => {
        const url =`${window.location.origin}/api`
        setBaseUrl(url)
        setSessionToken(getSessionToken())
    }, [])

    const {data, error} = useSWR((baseUrl && sessionToken) ? `${baseUrl}${api.list}?id=${listId}` : null, fetcher)

    if (error) return <div>Error</div>
    if (!data) {
        return <div>Loading</div>
    } else if (data.message !== "Successfully fetched users!") return <div>Error</div>

    return (
        <>
            <h1 className="text-2xl font-bold px-2 py-5">{listName}</h1>

            <hr />
            <br />

            <ModeratorTable listId={listId} listName={listName} modsDetails={moderators} modsList={modList} />

            <br />

            <ModUserTable listId={listId} listName={listName} userList={data.data} modsList={modList} />

        </>
    )
}