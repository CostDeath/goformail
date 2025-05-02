
import {expect, test, vitest} from "vitest";
import {render, screen} from "@testing-library/react";
import EmailApprovalForm from "@/components/emailApprovalRequests/emailApprovalForm";
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

test("Email view is rendered", () => {
    render(<EmailApprovalForm />);

    EmailViewChecker()
    expect(screen.getByRole("button", {name: "Approve"})).toBeDefined();
})