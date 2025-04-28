import {expect, test} from "vitest";
import {render, screen} from "@testing-library/react";
import Page from "@/app/(dashboards)/management/page";

test("Management Page is rendered", () => {
    render(<Page />);

    expect(screen.getByRole("link", {name: "+ New User"})).toBeDefined();
    expect(screen.getByTestId("table-head")).toBeDefined();
    expect(screen.getByTestId("table-body")).toBeDefined();
})