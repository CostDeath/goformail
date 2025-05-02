"use client"


export default function EmailView() {

    return (
        <>
            <div className="">
                <div data-testid="soc-email" className="text-xl py-3"> {/* Title of email goes here */}</div>
                <div data-testid="email-title" className="border-b border-black py-2">To: mailingList@example.com</div>
                <div data-testid="email-subject" className="border-b border-black py-2">Subject:</div>
            </div>
            <div data-testid="email-content" className="border border-black py-5 px-4 min-h-[50vh] bg-neutral-700">
                Content here
            </div>
        </>
    )
}