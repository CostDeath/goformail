"use client"

import {PageName, switchPage} from "@/states/linkStateHandler";
import Navbar from "@/components/navbar";

export default function Page() {
    const currentComponent = switchPage((state) => state.component)
    const currentPageName = switchPage((state) => state.name)
    //TODO: add additional condition bypassing this conditional if user opted to stay signed in
    if (currentPageName === PageName.LOGINSIGNUP) {
        return currentComponent
    }
    return (
        <div className="flex h-screen flex-col md:flex-row md:overflow-hidden">
            <div className="w-full flex-none md:w-64">
                <Navbar/>
            </div>
            <div className="flex-grow p-6 md:overflow-y-auto md:p-12">{currentComponent}</div>
        </div>
)
}