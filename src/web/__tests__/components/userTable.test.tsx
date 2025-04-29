import {expect, test, vitest} from "vitest";
import {render} from "@testing-library/react";
import UserTable from "@/components/management/userTable";
import useSWR from "swr";

vitest.mock("swr")

test("User table is loading", async () => {
    useSWR.mockReturnValue({
        data: undefined
    })
    const wrapper = render(<UserTable />);

    expect(wrapper.getByText("Loading")).toBeDefined();
    wrapper.unmount()

})

test("User table has loaded", async () => {
    useSWR.mockReturnValue({
        data: {message: "Successfully fetched users!", data: [{id: 1, email: "test@user.com"}]}
    })

    const wrapper = render(<UserTable />)

    vitest.useFakeTimers()
    vitest.runAllTimers()
    expect(wrapper.getByTestId("table-head")).toBeDefined();
    expect(wrapper.getByTestId("table-body")).toBeDefined();
    expect(wrapper.getByRole("link", {name: "test@user.com"})).toBeDefined()
    wrapper.unmount()
})

test("User table has loaded but given data was invalid", async () => {
    useSWR.mockReturnValue({
        data: {message: "Error"}
    })

    const wrapper = render(<UserTable />)

    vitest.useFakeTimers()
    vitest.runAllTimers()

    expect(wrapper.getByText("Error")).toBeDefined();
    wrapper.unmount()
})

test("User table has loaded but there is no data to show", async () => {
    useSWR.mockReturnValue({
        data: {message: "Successfully fetched users!"}
    })

    const wrapper = render(<UserTable />)

    vitest.useFakeTimers()
    vitest.runAllTimers()
    expect(wrapper.getByTestId("table-head")).toBeDefined();
    expect(wrapper.getByTestId("table-body")).toBeDefined();
    expect(wrapper.getByText("No Data to Show")).toBeDefined();
    wrapper.unmount()
})
