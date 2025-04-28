import {expect, test} from "vitest";
import {render, screen} from "@testing-library/react";
import UserTable from "@/components/management/userTable";

test("User table is rendered", () => {
    render(<UserTable />);

    expect(screen.getByTestId("table-head")).toBeDefined();
    expect(screen.getByTestId("table-body")).toBeDefined();
})
