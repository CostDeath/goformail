import {expect, test, vitest} from "vitest";
import {fireEvent, render, screen} from "@testing-library/react";
import Page from "@/app/(dashboards)/approvals/page";


vitest.mock("next/navigation", () => {
    const actual = vitest.importActual("next/navigation");
    return {
        ...actual,
        useRouter: vitest.fn(() => ({
            push: vitest.fn(),
        })),
        useSearchParams: vitest.fn(() => ({
            get: vitest.fn()
        })),
        usePathname: vitest.fn(() => ({
            usePathname: vitest.fn()
        }))
    }
})

test("Approval Requests page is rendered correctly and modal can be toggled", async () => {
    render(<Page />);

    expect(screen.queryByTestId("modal")).toBeNull();
    fireEvent.click(screen.getByText("exampleentry@email.com"));
    expect(screen.getByTestId("modal")).toBeDefined();
})
