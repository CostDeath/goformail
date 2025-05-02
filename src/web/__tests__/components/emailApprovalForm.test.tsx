
import {expect, test, vitest} from "vitest";
import {render} from "@testing-library/react";
import useSWR from "swr";
import EmailApprovalForm from "@/components/emailApprovalRequests/emailApprovalForm";

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

test("Email approval is loading", () => {
    useSWR.mockReturnValue({
        data: undefined
    })
    const wrapper = render(<EmailApprovalForm />);

    vitest.useFakeTimers()
    vitest.useFakeTimers()


    expect(wrapper.getByText("Loading")).toBeDefined();

    wrapper.unmount()
})

test("Email approval has rendered", () => {
    useSWR.mockReturnValue({
        data: {message: "Successfully fetched email!", data: {sender: "x@domain.tld", rcpt: ["y@domain.tld"], content: "content"}}
    })
    const wrapper = render(<EmailApprovalForm />);

    vitest.useFakeTimers()
    vitest.useFakeTimers()


    expect(wrapper.getByTestId("email-title")).toBeDefined();
    expect(wrapper.getByTestId("email-subject")).toBeDefined();
    expect(wrapper.getByTestId("email-content")).toBeDefined();

    wrapper.unmount()
})

test("Email approval has loaded but given data was invalid", () => {
    useSWR.mockReturnValue({
        data: {}
    })
    const wrapper = render(<EmailApprovalForm />);

    vitest.useFakeTimers()
    vitest.useFakeTimers()


    expect(wrapper.getByText("Error")).toBeDefined();

    wrapper.unmount()
})
