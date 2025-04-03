import {expect, test} from "vitest";
import {render, screen} from "@testing-library/react";
import Page from "@/app/page";

test("Whole page is rendered", () => {
    render(<Page />);
    expect(screen.getByRole("textbox", {name: "Email"})).toBeDefined();
    expect(screen.getByLabelText("Password")).toBeDefined();
    expect(screen.getByRole("button", {name: "Sign up!"})).toBeDefined();

    /*
    fireEvent.click(screen.getByRole("button", { name: "Sign up!" }));
    expect(screen.getByRole("textbox", {name: "Email"})).toBeDefined();
    expect(screen.getByRole("textbox", {name: "Student ID"})).toBeDefined();
    expect(screen.getByLabelText("Password")).toBeDefined();

     */
})
