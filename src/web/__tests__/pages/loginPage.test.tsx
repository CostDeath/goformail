import {expect, test} from "vitest";
import {render, screen} from "@testing-library/react";
import Page from "@/app/page";

test("Whole page is rendered", () => {
    render(<Page />);
    const emailBox = screen.getByRole("textbox", {name: "Email"});
    expect(emailBox).toBeDefined();
    expect(emailBox.getAttribute("type")).toBe("email");
    expect(emailBox.getAttribute("required")).toBeDefined();

    const passwordBox = screen.getByLabelText("Password");
    expect(passwordBox).toBeDefined();
    expect(passwordBox.getAttribute("type")).toBe("password");
    expect(passwordBox.getAttribute("required")).toBeDefined();


    expect(screen.getByTestId("to-sign-up")).toBeDefined();
})
