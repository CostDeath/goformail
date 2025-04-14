import {test, expect} from "vitest";
import {render} from "@testing-library/react";
import StyledInput from "@/components/styledInput";

test("Required styled input has been rendered", () => {
    const requiredWrapper = render(<StyledInput required />)
    const textBox = requiredWrapper.getByRole("textbox")
    expect(textBox).toBeDefined();
    expect(textBox.getAttribute("required")).toBeDefined();
    requiredWrapper.unmount();
})

test("Regular styled input has been rendered", () => {
    const regularWrapper = render(<StyledInput required={false} />)
    const textBox = regularWrapper.getByRole("textbox")
    expect(textBox).toBeDefined();
    expect(textBox.getAttribute("required")).toBeNull();
    regularWrapper.unmount();
})