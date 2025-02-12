import {render, screen} from "@testing-library/react";
import {expect, test} from "vitest";
import Card from "@/components/card";

test("Card is rendered", () => {
    render(<Card><h1>Mock component</h1></Card>);
    expect(screen.getByTestId("card")).toBeDefined();
    expect(screen.getByRole("heading", {level: 1, name: "Mock component"})).toBeDefined();
})