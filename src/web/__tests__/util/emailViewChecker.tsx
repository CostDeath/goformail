import {expect} from "vitest";
import {screen} from "@testing-library/react";

export const EmailViewChecker = () => {
    expect(screen.getByTestId("soc-email")).toBeDefined();
    expect(screen.getByTestId("email-title")).toBeDefined();
    expect(screen.getByTestId("email-subject")).toBeDefined();
    expect(screen.getByTestId("email-content")).toBeDefined();
}