import {expect, test} from "vitest";
import {render, screen} from "@testing-library/react";
import LoginForm from "@/components/loginSignup/loginForm";

test("Login form is rendered", () => {
    render(<LoginForm />);
    expect(screen.getByRole("textbox", {name: "Email"})).toBeDefined();
    expect(screen.getByLabelText("Password")).toBeDefined();
})