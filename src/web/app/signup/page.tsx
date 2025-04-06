import Card from "@/components/card";
import SignUpForm from "@/components/loginSignup/signUpForm";
import Link from "next/link";

export default function Page() {
    return (
        <div className="grid grid-rows-auto items-center justify-items-center min-h-screen p-8 pb-20 gap-16 sm:p-20">
            <div className="row-start-1">
                <Card><SignUpForm/></Card>
                <div className="grid grid-cols-2 py-5 px-5">
                    <p>Already have an account?</p>
                    <button className="text-blue-500 underline font-bold">
                        <Link href="/">Log in!</Link>
                    </button>
                </div>
            </div>
        </div>
    )
}