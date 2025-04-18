
import {expect, test} from "vitest";
import {render, screen} from "@testing-library/react";
import EmailApprovalForm from "@/components/emailApprovalRequests/emailApprovalForm";

test("Email view is rendered", () => {
    render(<EmailApprovalForm id={"1"} />);

    expect(screen.getByTestId("soc-email")).toBeDefined();
    expect(screen.getByTestId("email-title")).toBeDefined();
    expect(screen.getByTestId("email-subject")).toBeDefined();
    expect(screen.getByTestId("email-content")).toBeDefined();
    expect(screen.getByRole("button", {name: "Approve"})).toBeDefined();
    expect(screen.getByRole("button", {name: "Reject"})).toBeDefined();
})