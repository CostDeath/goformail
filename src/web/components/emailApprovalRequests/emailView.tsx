export default function EmailView() {
    // do fetches here
    return (
        <>
            <div className="py-6">
                <div data-testid="soc-email" className="text-xl py-6">Society email here</div>
                <div data-testid="email-title" className="border-b border-black py-2">Title:</div>
                <div data-testid="email-subject" className="border-b border-black py-2">Subject:</div>
            </div>
            <div data-testid="email-content" className="border border-black py-5 px-4 min-h-[50vh]">Content here</div>
            <div className="flex flex-row justify-end py-6 px-3">
                <div className="p-2">
                    <button className="bg-green-600 text-white p-2 rounded-xl hover:bg-green-700">Approve</button>
                </div>
                <div className="p-2">
                    <button className="bg-red-600 text-white p-2 px-4 rounded-xl hover:bg-red-700">Reject</button>
                </div>
            </div>
        </>
    )
}