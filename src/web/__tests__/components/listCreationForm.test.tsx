import {test, expect} from "vitest"
import {fireEvent, render, screen} from "@testing-library/react";
import ListCreationForm from "@/components/createList/listCreationForm";
import {EmailChecker} from "@/__tests__/util/formCheckers";

test("List Creation Form is rendered", () => {
    render(<ListCreationForm/>)

    EmailChecker("Mailing List Email")

    const noRecipientCheck = (name: string) => {
        expect(screen.queryByRole("textbox", {name: name})).toBeNull();
    }

    EmailChecker("recipient0")
    noRecipientCheck("recipient1")
    const delete1 = screen.getByRole("button", {name: /delete0/i})
    expect(delete1).toBeDefined();

    fireEvent.click(delete1)
    noRecipientCheck("recipient0")

    const addRecipient = screen.getByRole("button", {name: "+ Add recipient"})
    expect(addRecipient).toBeDefined();
    fireEvent.click(addRecipient);
    EmailChecker("recipient0")

    expect(screen.getByRole("button", {name: "Submit"})).toBeDefined();
})