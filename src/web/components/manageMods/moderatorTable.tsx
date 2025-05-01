import {api} from "@/components/api";
import {getSessionToken} from "@/components/sessionToken";


export default function ModeratorTable({listId, listName, modsDetails, modsList}: {
    listId: string | null,
    listName: string,
    modsDetails: {id: number, email: string}[],
    modsList: number[]
}) {
    const moderators = modsList
    const removeModerator = async (id: number) => {
        const url = `${window.location.origin}/api${api.list}?id=${listId}`
        const sessionToken = getSessionToken()
        for (let i = 0; i < moderators.length; i++) {
            if (id === moderators[i]) {
                moderators.splice(i, 1)
                break;
            }
        }
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
            alert("Successfully removed moderator")
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
                        <div data-testid="table-head-moderator" className="table-row bg-neutral-800/45 text-neutral-300">
                            <div className="table-cell border-b border-black text-left px-4 py-5 font-bold sm:pl-6">
                                Moderator
                            </div>
                            <div className="table-cell border-b border-black text-left px-4 py-5 font-bold sm:pl-6">
                                Action
                            </div>
                        </div>
                    </div>
                    <div data-testid="table-body-moderator" className="table-row-group">

                        {modsDetails && (
                            <>
                                {modsDetails.map((moderator) => (
                                    <div key={moderator.id}
                                         className="table-row shadow-inner text-neutral-300">
                                        <div className="table-cell black border-black py-3 text-sm">
                                            <div className="whitespace-nowrap py-3 pl-6 pr-3 flex items-center gap-3">
                                                {moderator.email}
                                            </div>
                                        </div>

                                        <div className="table-cell black border-black py-3 text-sm">
                                            <div className="whitespace-nowrap py-3 pl-6 pr-3 flex items-center gap-3">
                                                <button
                                                    className="bg-red-600 hover:bg-red-700 py-2 px-3 rounded-md"
                                                    onClick={() => removeModerator(moderator.id)}
                                                >
                                                    Remove
                                                </button>
                                            </div>
                                        </div>
                                    </div>
                                ))}
                            </>
                        )}

                        {!modsDetails && (
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