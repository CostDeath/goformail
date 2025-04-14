import {screen} from "@testing-library/react";
import {expect} from "vitest";

export function EmailChecker() {
    const emailBox = screen.getByRole("textbox", {name: "Email"});
    expect(emailBox).toBeDefined();
    expect(emailBox.getAttribute("type")).toBe("email");
    expect(emailBox.getAttribute("required")).toBeDefined();
}

export function PasswordChecker() {
    const passwordBox = screen.getByLabelText("Password");
    expect(passwordBox).toBeDefined();
    expect(passwordBox.getAttribute("type")).toBe("password");
    expect(passwordBox.getAttribute("required")).toBeDefined();
}