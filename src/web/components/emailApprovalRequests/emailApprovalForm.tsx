import EmailView from "@/components/emailView";
// import {useState} from "react";


export default function EmailApprovalForm() {
    /*
    //TODO: Remove this post example once doing fetching ticket
    const [title, setTitle] = useState("");
    const [body, setBody] = useState("");

    const submitApproval = async () => {
        let response = await fetch('https://jsonplaceholder.typicode.com/posts', {
            method: 'POST',
            body: JSON.stringify({
                title: title,
                body: body,
                userId: 1
            }),
            headers: {
                "Content-Type": "application/json"
            }
        })

        response = await response.json();
        alert(JSON.stringify(response))
    }

     */
    return (
        <>
            <EmailView />
            <div className="flex flex-row justify-end py-6 px-3">
                <div className="p-2">
                    <button className="bg-green-600 text-white p-2 rounded-xl hover:bg-green-700">Approve</button>
                </div>
                <div className="p-2">
                    <button className="bg-red-600 text-white p-2 px-4 rounded-xl hover:bg-red-700">Reject</button>
                </div>
            </div>


            {/*
            <br/>
            <br/>
            <h1>This is a post test</h1>
            <div>
                <h2>Testing posts</h2>
                <input type="text" value={title} onChange={e => setTitle(e.target.value)}/>
                <input type="text" value={body} onChange={e => setBody(e.target.value)}/>

                <button onClick={submitApproval}>Submit</button>
            </div>
            */}
        </>
    )
}