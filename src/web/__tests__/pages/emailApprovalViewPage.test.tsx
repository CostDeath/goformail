import {expect, test, vitest} from "vitest";
import {render, screen} from "@testing-library/react";
import Page from "@/app/(dashboards)/approvals/email/page";
import {EmailViewChecker} from "@/__tests__/util/emailViewChecker";

vitest.mock("next/navigation", () => {
    const actual = vitest.importActual("next/navigation");
    return {
        ...actual,
        useSearchParams: vitest.fn(() => ({
            get: (key: string) => {
                if (key === "id") return "1"
                return null
            }
        })),
    }
})

test("Email approval view page is rendered", () => {
    render(<Page />);

    EmailViewChecker()
    expect(screen.getByRole("button", {name: "Approve"})).toBeDefined();
})