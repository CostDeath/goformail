
import {expect, test} from "vitest";
import {render, screen} from "@testing-library/react";
import EmailApprovalForm from "@/components/emailApprovalRequests/emailApprovalForm";
import {EmailViewChecker} from "@/__tests__/util/emailViewChecker";

test("Email view is rendered", () => {
    render(<EmailApprovalForm id={"1"} />);

    EmailViewChecker()
    expect(screen.getByRole("button", {name: "Approve"})).toBeDefined();
    expect(screen.getByRole("button", {name: "Reject"})).toBeDefined();
})