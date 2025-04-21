import {expect, test} from "vitest";
import {render, screen} from "@testing-library/react";
import Page from "@/app/(dashboards)/approvals/email/page";

test("Email approval view page is rendered", () => {
    render(<Page />);

    expect(screen.getByTestId("soc-email")).toBeDefined();
    expect(screen.getByTestId("email-title")).toBeDefined();
    expect(screen.getByTestId("email-subject")).toBeDefined();
    expect(screen.getByTestId("email-content")).toBeDefined();
    expect(screen.getByRole("button", {name: "Approve"})).toBeDefined();
    expect(screen.getByRole("button", {name: "Reject"})).toBeDefined();
})