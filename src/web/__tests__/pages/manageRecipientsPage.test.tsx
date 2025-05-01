import {expect, test, vitest} from "vitest";
import {fireEvent, render, screen} from "@testing-library/react";
import useSWR from "swr";
import {EmailChecker} from "@/__tests__/util/formCheckers";
import Page from "@/app/(dashboards)/mailingLists/list/manageRecipients/page";

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

test("recipient form page is loading", async () => {
    useSWR.mockReturnValue({
        data: undefined
    })
    const wrapper = render(<Page/>);

    expect(wrapper.getByText("Loading")).toBeDefined();

    vitest.useFakeTimers()
    vitest.runAllTimers()

    wrapper.unmount()
})


test("recipient form page has loaded and is rendered", async () => {
    useSWR.mockReturnValue({
        data: {message: "Successfully fetched list!", data: {recipients: ["someEmail@example.com"]}}
    })
    const wrapper = render(<Page />)

    vitest.useFakeTimers()
    vitest.runAllTimers()

    const noRecipientCheck = (name: string) => {
        expect(screen.queryByRole("textbox", {name: name})).toBeNull();
    }

    noRecipientCheck("recipient1")
    const delete1 = screen.getByRole("button", {name: /delete0/i})
    expect(delete1).toBeDefined();

    fireEvent.click(delete1)
    noRecipientCheck("recipient0")

    const addRecipient = screen.getByRole("button", {name: "+ Add Another Recipient"})
    expect(addRecipient).toBeDefined();
    fireEvent.click(addRecipient);
    EmailChecker("recipient0")

    expect(wrapper.getByRole("button", {name: "Submit"})).toBeDefined();
    wrapper.unmount()
})

test("recipient form page has loaded but given data was invalid", async () => {
    useSWR.mockReturnValue({
        data: {}
    })
    const wrapper = render(<Page />)

    vitest.useFakeTimers()
    vitest.runAllTimers()
    expect(wrapper.getByText("Error")).toBeDefined();
    wrapper.unmount()
})