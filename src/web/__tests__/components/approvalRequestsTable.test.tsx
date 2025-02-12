import {expect, test, vitest} from "vitest";
import {render, screen} from "@testing-library/react";
import ApprovalRequestsTable from "@/components/emailApprovalRequests/approvalRequestsTable";


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

test("approval request table is rendered", () => {
    render(<ApprovalRequestsTable api="dummyString" />);

    expect(screen.getByTestId("table-head")).toBeDefined();
    expect(screen.getByTestId("table-body")).toBeDefined();
})