import {useSearchParams} from "next/navigation";
import {useModal} from "@/states/modalStateHandler";
import Modal from "@/components/modal";


export default function DeleteUser() {
    const search = useSearchParams()
    const id = search.get("id")
    const showModal = useModal((state) => state.toggled)
    const toggleModal = useModal((state) => state.toggleModal)
    return (
        <div className="px-3">
            <button className="bg-red-600 hover:bg-red-700 py-2 px-3 rounded-md"
            onClick={() => toggleModal(true)}>
                Delete User
            </button>

            {showModal && (
                <Modal width="150vh" height="85vh">
                    <h1 className="text-3xl py-5 text-center">Delete User?</h1>
                    <div className="p-5 text-xl">Are you sure you want to delete this user?</div>
                    <br/>
                    <footer className="p-4 flex flex-row justify-end">
                        <div className="py-2">
                            <button className="bg-red-600 hover:bg-red-700 py-2 px-3 rounded-md">Delete</button>
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