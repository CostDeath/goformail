import {expect, test} from "vitest";
import {render, screen} from "@testing-library/react";
import MailingListTable from "@/components/mailingLists/malingListTable";

test("mailing list table is rendered", () => {
    render(<MailingListTable api="dummyAPI" />);

    expect(screen.getByTestId("table-head")).toBeDefined();
    expect(screen.getByTestId("table-body")).toBeDefined();
})
