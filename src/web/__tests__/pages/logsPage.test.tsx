import {expect, test, vitest} from "vitest";
import {render, screen} from "@testing-library/react";
import useSWR from "swr";
import Page from "@/app/(dashboards)/management/logs/page"

test("Logs page is loading", async () => {
    useSWR.mockReturnValue({
        data: undefined
    })
    const wrapper = render(<Page/>);

    expect(wrapper.getByText("Loading")).toBeDefined();

    vitest.useFakeTimers()
    vitest.runAllTimers()

    wrapper.unmount()
})

vitest.mock("swr")

test("Logs page has loaded and is rendered", async () => {
    useSWR.mockReturnValue({
        data: [{id: "1"}]
    })
    const wrapper = render(<Page />)

    vitest.useFakeTimers()
    vitest.runAllTimers()
    expect(screen.getByTestId("table-head")).toBeDefined();
    expect(screen.getByTestId("table-body")).toBeDefined();
    wrapper.unmount()
})

test("Logs page has loaded but given data was invalid", async () => {
    useSWR.mockReturnValue({
        error: true
    })
    const wrapper = render(<Page />)

    vitest.useFakeTimers()
    vitest.runAllTimers()
    expect(wrapper.getByText("Error")).toBeDefined();
    wrapper.unmount()
})