import {expect, test, vitest} from "vitest";
import {render, screen} from "@testing-library/react";
import LogsTable from "@/components/logs/logsTable";
import useSWR from "swr";

test("Log table is loading", async () => {
    useSWR.mockReturnValue({
        data: undefined
    })
    const wrapper = render(<LogsTable/>);

    expect(wrapper.getByText("Loading")).toBeDefined();

    vitest.useFakeTimers()
    vitest.runAllTimers()

    wrapper.unmount()
})

vitest.mock("swr")

test("Log table component has loaded and is rendered", async () => {
    useSWR.mockReturnValue({
        data: [{id: "1"}]
    })
    const wrapper = render(<LogsTable />)

    vitest.useFakeTimers()
    vitest.runAllTimers()
    expect(screen.getByTestId("table-head")).toBeDefined();
    expect(screen.getByTestId("table-body")).toBeDefined();
    wrapper.unmount()
})

test("Log table component has loaded but given data was invalid", async () => {
    useSWR.mockReturnValue({
        error: true
    })
    const wrapper = render(<LogsTable />)

    vitest.useFakeTimers()
    vitest.runAllTimers()
    expect(wrapper.getByText("Error")).toBeDefined();
    wrapper.unmount()
})