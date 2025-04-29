import {useModal} from "@/states/modalStateHandler";
import Modal from "@/components/modal";
import {api} from "@/components/api";
import {redirect} from "next/navigation";
import {LinkTo} from "@/components/pageEnums";


export default function DeleteList({id}: {id: string | null}) {
    const showModal = useModal((state) => state.toggled)
    const toggleModal = useModal((state) => state.toggleModal)

    const deleteList = async () => {
        const response = await fetch(`${window.location.origin}/api${api.list}?id=${id}`, {
            method: "DELETE",
        })

        if (response.ok) {
            const result = await response.json();
            alert(result.message)
            redirect(LinkTo.MAILINGLISTS)
        } else {
            const result = await response.text();
            alert(result)
        }
    }



    return (
        <>
            <div>
                <button className="bg-red-600 hover:bg-red-700 py-2 px-3 rounded-md" onClick={() => toggleModal(true)}>
                    Delete Mailing List
                </button>
            </div>

            {showModal && id && (
                <Modal width="150vh" height="85vh">
                    <h1 data-testid="Delete Modal Header" className="text-3xl py-5 text-center">Delete Mailing List?</h1>
                    <div className="p-5 text-xl">Are you sure you want to delete this mailing list?</div>
                    <br/>
                    <footer className="p-4 flex flex-row justify-end">
                        <div className="p-2">
                            <button onClick={deleteList} className="bg-red-600 hover:bg-red-700 py-2 px-3 rounded-md">Delete</button>
                        </div>
                        <div className="p-2">
                            <button onClick={() => toggleModal(false)} className="border hover:bg-neutral-500 py-2 px-3 rounded-md">Cancel</button>
                        </div>

                    </footer>
                </Modal>
            )}

            {showModal && !id && (
                <Modal width="150vh" height="85vh">
                    <h1 data-testid="Delete Modal Header" className="text-3xl py-5 text-center">Error</h1>
                    <div className="p-5 text-xl">An error has occurred, the list may have already been deleted, please refresh the page</div>
                    <br/>
                    <footer className="p-4 flex flex-row justify-end">
                        <div className="p-2">
                            <button onClick={() => toggleModal(false)} className="border hover:bg-neutral-500 py-2 px-3 rounded-md">Cancel</button>
                        </div>

                    </footer>
                </Modal>
            )}
        </>
    )
}