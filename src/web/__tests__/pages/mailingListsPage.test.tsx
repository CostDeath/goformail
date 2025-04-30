import {expect, test, vitest} from "vitest";
import {render} from "@testing-library/react";
import MailingListsPage from "@/app/(dashboards)/mailingLists/page";
import useSWR from "swr";

vitest.mock("swr")

test("Mailing list page is loading", async () => {
    useSWR.mockReturnValue({
        data: undefined
    })
    const wrapper = render(<MailingListsPage />);

    expect(wrapper.getByText("Loading")).toBeDefined();

    vitest.useFakeTimers()
    vitest.runAllTimers()

    wrapper.unmount()
})



test("Mailing list page has loaded", async () => {
    useSWR.mockReturnValue({
        data: {message: "Successfully fetched lists!", data: [{id: 1, name: "test"}]}
    })

    const wrapper = render(<MailingListsPage />);

    vitest.useFakeTimers()
    vitest.runAllTimers()


    expect(wrapper.getByTestId("table-head")).toBeDefined();
    expect(wrapper.getByTestId("table-body")).toBeDefined();
    expect(wrapper.getByRole("link", {name: "test"})).toBeDefined();
    wrapper.unmount()
})

test("Mailing list page has loaded but given data was invalid", async () => {
    useSWR.mockReturnValue({
        data: {message: "Error"}
    })

    const wrapper = render(<MailingListsPage />);


    vitest.useFakeTimers()
    vitest.runAllTimers()


    expect(wrapper.getByText("Error")).toBeDefined();
    wrapper.unmount()
})

test("Mailing list page has loaded but no existing list data", async () => {
    useSWR.mockReturnValue({
        data: {message: "Successfully fetched lists!"}
    })

    const wrapper = render(<MailingListsPage />);


    vitest.useFakeTimers()
    vitest.runAllTimers()

    expect(wrapper.getByTestId("table-head")).toBeDefined();
    expect(wrapper.getByTestId("table-body")).toBeDefined();
    expect(wrapper.getByText("No Data to Show")).toBeDefined();
    wrapper.unmount()
})