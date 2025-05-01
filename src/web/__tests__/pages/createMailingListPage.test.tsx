import {test, expect} from "vitest"
import {fireEvent, render, screen} from "@testing-library/react";
import Page from "@/app/(dashboards)/mailingLists/create/page"
import {EmailChecker} from "@/__tests__/util/formCheckers";

test("List Creation page is rendered", () => {
    render(<Page/>)

    EmailChecker("Mailing List Name")

    const noRecipientCheck = (name: string) => {
        expect(screen.queryByRole("textbox", {name: name})).toBeNull();
    }

    EmailChecker("sender0")
    noRecipientCheck("sender1")
    const delete1 = screen.getByRole("button", {name: /delete0/i})
    expect(delete1).toBeDefined();

    fireEvent.click(delete1)
    noRecipientCheck("sender0")

    const addSender= screen.getByRole("button", {name: "+ Add Another Sender"})
    expect(addSender).toBeDefined();
    fireEvent.click(addSender);
    EmailChecker("sender0")

    expect(screen.getByRole("button", {name: "Submit"})).toBeDefined();
})