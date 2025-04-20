import {fireEvent, render, screen} from "@testing-library/react";
import {expect, test, vitest} from "vitest";
import DeleteList from "@/components/editList/deleteList";

vitest.mock("next/navigation", () => {
    const actual = vitest.importActual("next/navigation");
    return {
        ...actual,
        useSearchParams: vitest.fn(() => ({
            get: vitest.fn()
        })),
    }
})

test("Delete List is rendered", () => {
    render(<DeleteList />);

    const deleteButton = screen.getByRole("button", { name: "Delete Mailing List" })
    expect(deleteButton).toBeDefined();
    expect(screen.queryByTestId("Delete Modal Header")).toBeNull();
    fireEvent.click(deleteButton)
    expect(screen.getByTestId("Delete Modal Header")).toBeDefined();
})

