import {expect, test} from "vitest";
import {render, screen} from "@testing-library/react";
import SignUpForm from "@/components/loginSignup/signUpForm";

test("Sign up form is rendered", () => {
    render(<SignUpForm />);
    expect(screen.getByRole("textbox", {name: "Email"})).toBeDefined();
    expect(screen.getByRole("textbox", {name: "Student ID"})).toBeDefined();
    expect(screen.getByLabelText("Password")).toBeDefined();
})