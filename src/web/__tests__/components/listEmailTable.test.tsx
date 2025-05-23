import {expect, test, vitest} from "vitest";
import {render, screen} from "@testing-library/react";
import useSWR from "swr";
import ListEmailsTable from "@/components/list/listEmailsTable";

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


test("List Email Table is loading", async () => {
    useSWR.mockReturnValue({
        data: undefined
    })
    const wrapper = render(<ListEmailsTable/>);

    expect(wrapper.getByText("Loading...")).toBeDefined();

    vitest.useFakeTimers()
    vitest.runAllTimers()

    wrapper.unmount()
})

vitest.mock("swr")

test("List Email Table component has loaded and is rendered", async () => {
    useSWR.mockReturnValue({
        data: {message: "Successfully fetched emails!", data: {offset: 0, emails: []}}
    })
    const wrapper = render(<ListEmailsTable />)

    vitest.useFakeTimers()
    vitest.runAllTimers()

    expect(screen.getByTestId("table-head-emails")).toBeDefined();
    expect(screen.getByTestId("table-body-emails")).toBeDefined();
    wrapper.unmount()
})

test("List Email Table component has loaded but given data was invalid", async () => {
    useSWR.mockReturnValue({
        data: {}
    })
    const wrapper = render(<ListEmailsTable />)

    vitest.useFakeTimers()
    vitest.runAllTimers()
    expect(wrapper.getByText("Error loading emails list")).toBeDefined();
    wrapper.unmount()
})