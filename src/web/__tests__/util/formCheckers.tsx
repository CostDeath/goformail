import {screen} from "@testing-library/react";
import {expect} from "vitest";

export function EmailChecker(label?: string) {
    const email = (label) ? label : "Email"
    const emailBox = screen.getByRole("textbox", {name: label});
    expect(emailBox).toBeDefined();
    expect(emailBox.getAttribute("type")).toBe("email");
    expect(emailBox.getAttribute("required")).toBeDefined();
}

export function PasswordChecker(label?: string) {
    const password = (label) ? label : "Password"
    const passwordBox = screen.getByLabelText(password);
    expect(passwordBox).toBeDefined();
    expect(passwordBox.getAttribute("type")).toBe("password");
    expect(passwordBox.getAttribute("required")).toBeDefined();
}
