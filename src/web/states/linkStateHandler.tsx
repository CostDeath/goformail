import {create} from "zustand"
import {ReactNode} from "react";
import LoginSignupPage from "@/app/pages/LoginSignupPage";

interface LinkState {
    name: string;
    component: ReactNode;
    changePage: (page: ReactNode, pageName: string) => void;
}

//enums for names associated with each
export const PageName = {
    LOGINSIGNUP: "Login/Signup",
    MAILINGLISTS: "Mailing Lists",
    APPROVALREQUESTS: "Email Approval Requests",
}

export const switchPage = create<LinkState>()((set) => ({
    name: PageName.LOGINSIGNUP,
    component: <LoginSignupPage />,
    changePage: (page: ReactNode, pageName: string) => set((state) => ({component: page, name: pageName})),
}))