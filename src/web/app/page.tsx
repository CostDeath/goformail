import Card from "@/components/card";
import LoginForm from "@/components/loginSignup/loginForm";
import Link from "next/link";
import {LinkTo} from "@/states/linkStateHandler";

export default function Page() {
    return (
        <div className="grid grid-rows-auto items-center justify-items-center min-h-screen p-8 pb-20 gap-16 sm:p-20">
            <div className="row-start-1">
                <Card><LoginForm/></Card>
                <div className="grid grid-cols-2 py-5 px-5">
                    <p>Don&apos;t have an account?</p>
                    <button className="text-blue-500 underline font-bold">
                        <Link href={LinkTo.SIGNUP}>Sign up!</Link>
                    </button>
                </div>
            </div>
        </div>
    )
}
