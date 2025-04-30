import {redirect} from "next/navigation";
import {LinkTo} from "@/components/pageEnums";

export function handleInvalidSessionToken() {
    localStorage.removeItem("sessionToken")
    redirect(LinkTo.LOGIN)
}

export function getSessionToken() {
    return localStorage.getItem("sessionToken")
}