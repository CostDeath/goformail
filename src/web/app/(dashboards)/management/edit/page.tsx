import Card from "@/components/card";
import {Suspense} from "react";
import EditUserForm from "@/components/editUser/editUserForm";

export default function Page() {
    return (
        <>
            <Card>
                <Suspense>
                    <EditUserForm />
                </Suspense>
            </Card>
        </>
    )
}