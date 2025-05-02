import ApprovalRequestsTable from "@/components/emailApprovalRequests/approvalRequestsTable";
import {expect, test, vitest} from "vitest";
import {render} from "@testing-library/react";
import useSWR from "swr";

vitest.mock("swr")

test("Approval Requests table is Loading", async () => {
    useSWR.mockReturnValue({
        data: undefined
    })
    const wrapper = render(<ApprovalRequestsTable />);

    vitest.useFakeTimers()
    vitest.runAllTimers()

    expect(wrapper.getByText("Loading")).toBeDefined();

    wrapper.unmount()
})

test("Approval Requests table has rendered", async () => {
    useSWR.mockReturnValue({
        data: {message: "Successfully fetched emails!", data: {offset: 0, emails: []}}
    })
    const wrapper = render(<ApprovalRequestsTable />);


    expect(wrapper.getByTestId("table-head")).toBeDefined();
    expect(wrapper.getByTestId("table-body")).toBeDefined();

    wrapper.unmount()
})

test("Approval Requests table has loaded but given data was invalid", async () => {
    useSWR.mockReturnValue({
        data: {}
    })
    const wrapper = render(<ApprovalRequestsTable />);

    vitest.useFakeTimers()
    vitest.runAllTimers()

    expect(wrapper.getByText("Error")).toBeDefined();

    wrapper.unmount()
})
