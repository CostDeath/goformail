import {useModal} from "@/states/modalStateHandler";
import Modal from "@/components/modal";
import {useSearchParams} from "next/navigation";


export default function DeleteList() {
    const search = useSearchParams();
    const id = search.get("id");
    const showModal = useModal((state) => state.toggled)
    const toggleModal = useModal((state) => state.toggleModal)
    console.log(id)
    return (
        <>
            <div>
                <button className="bg-red-600 hover:bg-red-700 py-2 px-3 rounded-md" onClick={() => toggleModal(true)}>
                    Delete Mailing List
                </button>
            </div>

            {showModal && (
                <Modal width="150vh" height="85vh">
                    <h1 className="text-3xl py-5 text-center">Delete Mailing List?</h1>
                    <div className="p-5 text-xl">Are you sure you want to delete this mailing list?</div>
                    <br/>
                    <footer className="p-4 flex flex-row justify-end">
                        <div className="p-2">
                            <button className="bg-red-600 hover:bg-red-700 py-2 px-3 rounded-md">Delete</button>
                        </div>
                        <div className="p-2">
                            <button onClick={() => toggleModal(false)} className="border hover:bg-neutral-500 py-2 px-3 rounded-md">Cancel</button>
                        </div>

                    </footer>
                </Modal>
            )}
        </>
    )
}