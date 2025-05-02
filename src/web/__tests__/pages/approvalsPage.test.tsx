import {expect, test, vitest} from "vitest";
import {render} from "@testing-library/react";
import Page from "@/app/(dashboards)/approvals/page";
import useSWR from "swr";

vitest.mock("swr")

test("Approval Requests page is Loading", async () => {
    useSWR.mockReturnValue({
        data: undefined
    })
    const wrapper = render(<Page />);

    vitest.useFakeTimers()
    vitest.runAllTimers()

    expect(wrapper.getByText("Loading")).toBeDefined();

    wrapper.unmount()
})

test("Approval Requests page has rendered", async () => {
    useSWR.mockReturnValue({
        data: {message: "Successfully fetched emails!", data: {offset: 0, emails: []}}
    })
    const wrapper = render(<Page />);


    expect(wrapper.getByTestId("table-head")).toBeDefined();
    expect(wrapper.getByTestId("table-body")).toBeDefined();

    wrapper.unmount()
})

test("Approval Requests page has loaded but given data was invalid", async () => {
    useSWR.mockReturnValue({
        data: {}
    })
    const wrapper = render(<Page />);

    vitest.useFakeTimers()
    vitest.runAllTimers()

    expect(wrapper.getByText("Error")).toBeDefined();

    wrapper.unmount()
})
