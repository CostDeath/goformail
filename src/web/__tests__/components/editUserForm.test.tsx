import {expect, test, vitest} from "vitest";
import {render} from "@testing-library/react";
import EditUserForm from "@/components/editUser/editUserForm";
import useSWR from "swr";
import {EmailChecker} from "@/__tests__/util/formCheckers";
import {permissionsList} from "@/components/permissions";

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


test("Edit User is loading", async () => {
    useSWR.mockReturnValue({
        data: undefined
    })
    const wrapper = render(<EditUserForm/>);

    expect(wrapper.getByText("Loading")).toBeDefined();

    vitest.useFakeTimers()
    vitest.runAllTimers()

    wrapper.unmount()
})

vitest.mock("swr")

test("Edit User component has loaded and is rendered", async () => {
    useSWR.mockReturnValue({
        data: {id: "1"}
    })
    const wrapper = render(<EditUserForm />)

    vitest.useFakeTimers()
    vitest.runAllTimers()
    EmailChecker()
    expect(wrapper.getByRole("button", {name: "Delete User"})).toBeDefined();

    permissionsList.forEach(permission => {
        expect(wrapper.getByRole("checkbox", {name: permission.label})).toBeDefined();
    })

    expect(wrapper.getByRole("button", {name: "Edit User"})).toBeDefined();
    wrapper.unmount()
})

test("Edit User component has loaded but given data was invalid", async () => {
    useSWR.mockReturnValue({
        data: {}
    })
    const wrapper = render(<EditUserForm />)

    vitest.useFakeTimers()
    vitest.runAllTimers()
    expect(wrapper.getByText("Error")).toBeDefined();
    wrapper.unmount()
})