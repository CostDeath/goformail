import {expect, test, vitest} from "vitest";
import {render} from "@testing-library/react";
import useSWR from "swr";
import Page from "@/app/(dashboards)/approvals/email/page";

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

test("Email approval page is loading", () => {
    useSWR.mockReturnValue({
        data: undefined
    })
    const wrapper = render(<Page />);

    vitest.useFakeTimers()
    vitest.useFakeTimers()


    expect(wrapper.getByText("Loading")).toBeDefined();

    wrapper.unmount()
})

test("Email approval page has rendered", () => {
    useSWR.mockReturnValue({
        data: {message: "Successfully fetched email!", data: {sender: "x@domain.tld", rcpt: ["y@domain.tld"], content: "content"}}
    })
    const wrapper = render(<Page />);

    vitest.useFakeTimers()
    vitest.useFakeTimers()


    expect(wrapper.getByTestId("email-title")).toBeDefined();
    expect(wrapper.getByTestId("email-subject")).toBeDefined();
    expect(wrapper.getByTestId("email-content")).toBeDefined();

    wrapper.unmount()
})

test("Email approval page has loaded but given data was invalid", () => {
    useSWR.mockReturnValue({
        data: {}
    })
    const wrapper = render(<Page />);

    vitest.useFakeTimers()
    vitest.useFakeTimers()


    expect(wrapper.getByText("Error")).toBeDefined();

    wrapper.unmount()
})
