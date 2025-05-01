import {expect, test, vitest} from "vitest";
import {fireEvent, render, screen} from "@testing-library/react";
import useSWR from "swr";
import {EmailChecker} from "@/__tests__/util/formCheckers";
import ListEditForm from "@/components/editList/listEditForm";

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


test("Edit List is loading", async () => {
    useSWR.mockReturnValue({
        data: undefined
    })
    const wrapper = render(<ListEditForm/>);

    expect(wrapper.getByText("Loading")).toBeDefined();

    vitest.useFakeTimers()
    vitest.runAllTimers()

    wrapper.unmount()
})

vitest.mock("swr")

test("Edit List component has loaded and is rendered", async () => {
    useSWR.mockReturnValue({
        data: {message: "Successfully fetched list!", data: {name: "someEmail@example.com"}}
    })
    const wrapper = render(<ListEditForm />)

    vitest.useFakeTimers()
    vitest.runAllTimers()
    EmailChecker("Mailing List Name")
    expect(wrapper.getByRole("button", {name: "Delete Mailing List"})).toBeDefined();

    const noSenderCheck = (name: string) => {
        expect(screen.queryByRole("textbox", {name: name})).toBeNull();
    }

    noSenderCheck("sender1")
    const delete1 = screen.getByRole("button", {name: /delete0/i})
    expect(delete1).toBeDefined();

    fireEvent.click(delete1)
    noSenderCheck("sender0")

    const addSender = screen.getByRole("button", {name: "+ Add Another Sender"})
    expect(addSender).toBeDefined();
    fireEvent.click(addSender);
    EmailChecker("sender0")

    expect(wrapper.getByRole("button", {name: "Submit"})).toBeDefined();
    wrapper.unmount()
})

test("Edit List component has loaded but given data was invalid", async () => {
    useSWR.mockReturnValue({
        data: {}
    })
    const wrapper = render(<ListEditForm />)

    vitest.useFakeTimers()
    vitest.runAllTimers()
    expect(wrapper.getByText("Error")).toBeDefined();
    wrapper.unmount()
})