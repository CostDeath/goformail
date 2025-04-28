import {render, screen} from "@testing-library/react";
import {expect, test} from "vitest";
import Navbar from "@/components/navbar";
import {PageName} from "@/components/pageEnums";

test("navbar is rendered", () => {
    render(<Navbar />);
    const management = screen.getByTestId(PageName.MANAGEMENT)
    const mailingLists = screen.getByTestId(PageName.MAILINGLISTS)
    const emailApprovals = screen.getByTestId(PageName.APPROVALREQUESTS)
    const signOut = screen.getByTestId(PageName.LOGIN)

    expect(management).toBeDefined();
    expect(mailingLists).toBeDefined();
    expect(emailApprovals).toBeDefined();
    expect(signOut).toBeDefined();

    // May remove pagination so might remove these tests
    expect(mailingLists.getAttribute("onClick")).toBeDefined();
    expect(emailApprovals.getAttribute("onClick")).toBeDefined();
    expect(signOut.getAttribute("onClick")).toBeDefined();
})