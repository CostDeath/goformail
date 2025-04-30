import {expect, test, vitest} from "vitest";
import {render} from "@testing-library/react";
import MailingListTable from "@/components/mailingLists/malingListTable";
import useSWR from "swr";

vitest.mock("swr")

test("mailing list table is loading", async () => {
    useSWR.mockReturnValue({
        data: undefined
    })

    const wrapper = render(<MailingListTable />);

    expect(wrapper.getByText("Loading")).toBeDefined()

    wrapper.unmount()
})

test("mailing list table has loaded", async () => {
    useSWR.mockReturnValue({
        data: {message: "Successfully fetched lists!", data: [{id: 1, name: "test"}]}
    })

    const wrapper = render(<MailingListTable />)

    vitest.useFakeTimers()
    vitest.runAllTimers()

    expect(wrapper.getByTestId("table-head")).toBeDefined();
    expect(wrapper.getByTestId("table-body")).toBeDefined();
    expect(wrapper.getByRole("link", {name: "test"})).toBeDefined();
    wrapper.unmount()
})

test("mailing list table has loaded but given data was invalid", async () => {
    useSWR.mockReturnValue({
        data: {message: "Error"}
    })

    const wrapper = render(<MailingListTable />)

    vitest.useFakeTimers()
    vitest.runAllTimers()

    expect(wrapper.getByText("Error")).toBeDefined();
    wrapper.unmount()
})

test("mailing list table has loaded but no existing list data", async () => {
    useSWR.mockReturnValue({
        data: {message: "Successfully fetched lists!"}
    })

    const wrapper = render(<MailingListTable />)

    vitest.useFakeTimers()
    vitest.runAllTimers()

    expect(wrapper.getByTestId("table-head")).toBeDefined();
    expect(wrapper.getByTestId("table-body")).toBeDefined();
    expect(wrapper.getByText("No Data to Show")).toBeDefined();
    wrapper.unmount()
})