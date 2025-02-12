import {expect, test, vitest} from "vitest";
import {render, screen} from "@testing-library/react";
import MailingListsPage from "@/app/pages/mailingListsPage";


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

test("Mailing list page is rendered", async () => {
    render(MailingListsPage());

    expect(screen.getByTestId("table-head")).toBeDefined();
    expect(screen.getByTestId("table-body")).toBeDefined();
    expect(screen.getByTestId("pagination")).toBeDefined();
})