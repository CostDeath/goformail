import {expect, test, vitest} from "vitest";
import {render, screen} from "@testing-library/react";
import useSWR from "swr";
import Page from "@/app/(dashboards)/mailingLists/list/page";

vitest.mock("next/navigation", () => {
    const actual = vitest.importActual("next/navigation");
    return {
        ...actual,
        useSearchParams: vitest.fn(() => ({
            get: (key: string) => {
                if (key === "id") return "1"
                return null
            }
        })),
    }
})


test("List Email page is loading", async () => {
    useSWR.mockReturnValue({
        data: undefined
    })
    const wrapper = render(<Page/>);

    expect(wrapper.getByText("Loading...")).toBeDefined();

    vitest.useFakeTimers()
    vitest.runAllTimers()

    wrapper.unmount()
})

vitest.mock("swr")

test("List Email page has loaded and is rendered", async () => {
    useSWR.mockReturnValue({
        data: {id: "1"}
    })
    const wrapper = render(<Page />)

    vitest.useFakeTimers()
    vitest.runAllTimers()

    expect(screen.getByTestId("table-head")).toBeDefined();
    expect(screen.getByTestId("table-body")).toBeDefined();
    wrapper.unmount()
})

test("List Email page has loaded but given data was invalid", async () => {
    useSWR.mockReturnValue({
        data: {}
    })
    const wrapper = render(<Page />)

    vitest.useFakeTimers()
    vitest.runAllTimers()
    expect(wrapper.getByText("Error")).toBeDefined();
    wrapper.unmount()
})