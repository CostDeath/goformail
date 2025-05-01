import {expect, test} from "vitest";
import {render} from "@testing-library/react";
import ModUserTable from "@/components/manageMods/modUserTable";
import {User} from "@/models/user";


test("Mod User Table is loading", async () => {
    const user: User = {id: 1, email: "hello@example.com", permissions: []}
    const wrapper = render(<ModUserTable listId={"1"} listName={"example"} userList={[user]} modsList={[]}  />);

    expect(wrapper.getByTestId("table-head-user")).toBeDefined()
    expect(wrapper.getByTestId("table-body-user")).toBeDefined()
    expect(wrapper.getByRole("button", {name: "+ Add Moderator"})).toBeDefined()

    wrapper.unmount()
})
