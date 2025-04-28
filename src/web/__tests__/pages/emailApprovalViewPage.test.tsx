import {expect, test} from "vitest";
import {render, screen} from "@testing-library/react";
import Page from "@/app/(dashboards)/approvals/email/page";
import {EmailViewChecker} from "@/__tests__/util/emailViewChecker";

test("Email approval view page is rendered", () => {
    render(<Page />);

    EmailViewChecker()
    expect(screen.getByRole("button", {name: "Approve"})).toBeDefined();
    expect(screen.getByRole("button", {name: "Reject"})).toBeDefined();
})