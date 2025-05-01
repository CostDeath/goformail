import {expect, test, vitest} from "vitest";
import {render} from "@testing-library/react";
import useSWR from "swr";
import ManageModsForm from "@/components/manageMods/manageModsForm";

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

vitest.mock("swr")

test("Manage mods form is loading", async () => {
    useSWR.mockReturnValue({
        data: undefined
    })
    const wrapper = render(<ManageModsForm />);

    expect(wrapper.getByText("Loading")).toBeDefined();

    vitest.useFakeTimers()
    vitest.runAllTimers()

    wrapper.unmount()
})


test("Manage mods form component has loaded and is rendered", async () => {
    useSWR.mockReturnValue({
        data: {message: "Successfully fetched users!", data: [{id: 1, email: "hello@gmail.com", permissions: []}]}
    })
    const wrapper = render(<ManageModsForm />)

    vitest.useFakeTimers()
    vitest.runAllTimers()

    expect(wrapper.getByTestId("table-head-moderator")).toBeDefined();
    expect(wrapper.getByTestId("table-body-moderator")).toBeDefined();
    expect(wrapper.getByTestId("table-head-user")).toBeDefined();
    expect(wrapper.getByTestId("table-body-user")).toBeDefined();
    wrapper.unmount()
})

test("Manage Mods Form component has loaded but given data was invalid", async () => {
    useSWR.mockReturnValue({
        data: {}
    })
    const wrapper = render(<ManageModsForm />)

    vitest.useFakeTimers()
    vitest.runAllTimers()
    expect(wrapper.getByText("Error")).toBeDefined();
    wrapper.unmount()
})