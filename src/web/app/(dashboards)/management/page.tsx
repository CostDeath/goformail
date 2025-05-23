import UserTable from "@/components/management/userTable";


export default function Page() {
    return (
        <>
            <div className="py-5">
                <span>
                    <a href="/ui/management/add.html" className="bg-cyan-600 text-white hover:bg-cyan-500 px-3 py-2 rounded-md font-bold">+ New User</a>
                </span>
            </div>
            <UserTable />
        </>
    )
}