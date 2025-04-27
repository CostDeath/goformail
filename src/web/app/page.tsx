import Card from "@/components/card";
import LoginForm from "@/components/login/loginForm";

export default function Page() {
    return (
        <div className="grid grid-rows-auto items-center justify-items-center min-h-screen p-8 pb-20 gap-16 sm:p-20">
            <div className="row-start-1">
                <Card><LoginForm/></Card>
            </div>
        </div>
    )
}
