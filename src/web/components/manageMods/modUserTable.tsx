"use client"

import {User} from "@/models/user";
import {api} from "@/components/api";
import {getSessionToken} from "@/components/sessionToken";

export default function ModUserTable({listId, listName, userList, modsList}: {
    listId: string | null,
    listName: string,
    userList: User[],
    modsList: number[]
}) {
    const moderators = modsList
    const addModerator = async (id: number) => {
        const url = `${window.location.origin}/api${api.list}?id=${listId}`
        const sessionToken = getSessionToken()
        for (let i = 0; i < moderators.length; i++) {
            if (moderators[i] === id) {
                alert("Moderator already exists")
                return
            }
        }
        moderators.push(id)
        const response = await fetch(url, {
            method: "PATCH",
            body: JSON.stringify({
                name: listName,
                mods: moderators
            }),
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${sessionToken}`
            }
        })

        if (response.ok) {
            alert("Successfully added moderator")
        } else {
            const result = await response.text()
            alert(result)
        }
    }

    return (
        <>

            <div className="overflow-auto max-h-[75vh] shadow-[0_3px_10px_-1px_rgba(0,0,0,1)]">
                <div className="min-w-full table text-gray-900">
                    <div className="table-header-group text-left text-sm font-normal">
                        <div data-testid="table-head-user" className="table-row bg-neutral-800/45 text-neutral-300">
                            <div className="table-cell border-b border-black text-left px-4 py-5 font-bold sm:pl-6">
                                User
                            </div>
                            <div className="table-cell border-b border-black text-left px-4 py-5 font-bold sm:pl-6">
                                Action
                            </div>
                        </div>
                    </div>
                    <div data-testid="table-body-user" className="table-row-group">

                        {userList && (
                            <>
                                {userList.map((user) => (
                                    <div key={user.id}
                                         className="table-row shadow-inner text-neutral-300">
                                        <div className="table-cell black border-black py-3 text-sm">
                                            <div className="whitespace-nowrap py-3 pl-6 pr-3 flex items-center gap-3">
                                                {user.email}
                                            </div>
                                        </div>

                                        <div className="table-cell black border-black py-3 text-sm">
                                            <div className="whitespace-nowrap py-3 pl-6 pr-3 flex items-center gap-3">
                                                <button
                                                    className="bg-cyan-600 text-white hover:bg-cyan-500 py-2 px-3 rounded-md"
                                                    onClick={() => addModerator(user.id)}
                                                >
                                                    + Add Moderator
                                                </button>
                                            </div>
                                        </div>
                                    </div>
                                ))}
                            </>
                        )}

                        {!userList && (
                            <div className="table-row shadow-inner text-neutral-300">
                                <div className="table-cell border-black border-b py-3 text-sm">
                                    <div className="whitespace-nowrap py-3 pl-6 pr-3 flex items-center gap-3">
                                        No Data to Show
                                    </div>
                                </div>
                            </div>
                        )}
                    </div>
                </div>
            </div>
        </>
    )
}