import {useModal} from "@/states/modalStateHandler";
import Modal from "@/components/modal";
import {api} from "@/components/api";
import {redirect} from "next/navigation";
import {LinkTo} from "@/components/pageEnums";
import {getSessionToken} from "@/components/sessionToken";


export default function DeleteUser({id}: {id: number | null}) {
    const showModal = useModal((state) => state.toggled)
    const toggleModal = useModal((state) => state.toggleModal)

    const deleteUser = async () => {
        const sessionToken = getSessionToken();
        const response = await fetch(`${window.location.origin}/api${api.user}?id=${id}`, {
            method: "DELETE",
            headers: {
                "Authorization": `Bearer ${sessionToken}`
            }
        })

        if (response.ok) {
            const result = await response.json()
            alert(result.message)
            redirect(LinkTo.MANAGEMENT)
        } else {
            const result = await response.text()
            alert(result)
        }
    }
    return (
        <div className="px-3">
            <button className="bg-red-600 hover:bg-red-700 py-2 px-3 rounded-md"
            onClick={() => toggleModal(true)}>
                Delete User
            </button>

            {showModal && id && (
                <Modal width="150vh" height="85vh">
                    <h1 data-testid="Delete Modal Header" className="text-3xl py-5 text-center">Delete User?</h1>
                    <div className="p-5 text-xl">Are you sure you want to delete this user?</div>
                    <br/>
                    <footer className="p-4 flex flex-row justify-end">
                        <div className="py-2">
                            <button onClick={deleteUser} className="bg-red-600 hover:bg-red-700 py-2 px-3 rounded-md">Delete</button>
                        </div>
                        <div className="p-2">
                            <button onClick={() => toggleModal(false)}
                                    className="border hover:bg-neutral-500 py-2 px-3 rounded-md">Cancel
                            </button>
                        </div>

                    </footer>
                </Modal>
            )}
        </div>
    )
}