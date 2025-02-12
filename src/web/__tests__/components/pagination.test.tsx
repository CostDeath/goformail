import {vitest} from "vitest";
import {expect, test} from "vitest";
import {render, screen} from "@testing-library/react";
import Pagination from "@/components/pagination";

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

test("Pagination is rendered", () => {
    render(<Pagination totalPages={2}/>);

    expect(screen.getByTestId("pagination")).toBeDefined();
})