import {expect, test} from "vitest";
import {render} from "@testing-library/react";
import ModeratorTable from "@/components/manageMods/moderatorTable";

test("Moderator Table is loading", async () => {
    const wrapper = render(<ModeratorTable listId={"1"} listName={"example"} modsDetails={[{id: 1, email: "hello@example.com"}]} modsList={[1]} />);

    expect(wrapper.getByTestId("table-head-moderator")).toBeDefined()
    expect(wrapper.getByTestId("table-body-moderator")).toBeDefined()
    expect(wrapper.getByRole("button", {name: "Remove"})).toBeDefined()

    wrapper.unmount()
})