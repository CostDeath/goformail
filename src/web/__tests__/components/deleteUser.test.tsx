import {expect, test, vitest} from "vitest";
import {fireEvent, render, screen} from "@testing-library/react";
import DeleteUser from "@/components/editUser/deleteUser";

vitest.mock("next/navigation", () => {
    const actual = vitest.importActual("next/navigation");
    return {
        ...actual,
        useSearchParams: vitest.fn(() => ({
            get: vitest.fn()
        })),
    }
})

test("Delete User is rendered", () => {
    render(<DeleteUser/>);

    const deleteButton = screen.getByRole("button", { name: "Delete User" })
    expect(deleteButton).toBeDefined();
    expect(screen.queryByTestId("Delete Modal Header")).toBeNull();
    fireEvent.click(deleteButton)
    expect(screen.getByTestId("Delete Modal Header")).toBeDefined();
})