import Card from "@/components/card";
import LoginForm from "@/components/login/loginForm";
import Link from "next/link";
import {LinkTo} from "@/components/pageEnums";

export default function Page() {
    return (
        <div className="grid grid-rows-auto items-center justify-items-center min-h-screen p-8 pb-20 gap-16 sm:p-20">
            <div className="row-start-1">
                <Card><LoginForm/></Card>
                <div className="grid grid-cols-2 py-5 px-5">
                    <p>Don&apos;t have an account?</p>
                    <Link href={LinkTo.SIGNUP} className="text-blue-400 font-bold px-7" data-testid="to-sign-up">
                        Sign up!
                    </Link>
                </div>
            </div>
        </div>
    )
}
