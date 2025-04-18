
export default function EmailView({id}: {id: string}) {
    // fetch email content here
    return (
        <>
            <div className="py-6">
                <div data-testid="soc-email" className="text-xl py-6">Society email here</div>
                <div data-testid="email-title" className="border-b border-black py-2">Title: {id}</div>
                <div data-testid="email-subject" className="border-b border-black py-2">Subject:</div>
            </div>
            <div data-testid="email-content" className="border border-black py-5 px-4 min-h-[50vh]">Content here</div>
        </>
    )
}