import {expect, test} from "vitest";
import {render, screen} from "@testing-library/react";
import Modal from "@/components/modal";
import {vitest} from "vitest";

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

test("Modal is rendered", () => {
    render(<Modal width="100vh" height="90vh"><h1>Mock component</h1></Modal>);

    expect(screen.getByRole("button", { name: "X" })).toBeDefined();
    expect(screen.getByRole("heading", {level: 1, name: "Mock component"})).toBeDefined();
})