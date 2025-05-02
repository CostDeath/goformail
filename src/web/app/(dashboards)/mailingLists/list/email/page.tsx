import EmailView from "@/components/emailView";
import {Suspense} from "react";

export default function Page() {
    return (
        <>
            <Suspense>
                <EmailView />
            </Suspense>
        </>
    )
}