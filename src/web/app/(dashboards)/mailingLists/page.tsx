import MailingListTable from "@/components/mailingLists/malingListTable";

export default function Page() {
    return (
        <div className="w-full">
            <div className="py-5">
            <a href="/ui/mailingLists/create.html" className="bg-cyan-600 text-white hover:bg-cyan-500 px-3 py-2 rounded-md font-bold">+ New Mailing List</a>
            </div>
            <MailingListTable />
        </div>
    )
}