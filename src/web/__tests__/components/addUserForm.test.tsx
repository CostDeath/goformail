import {test, expect} from "vitest"
import AddUserForm from "@/components/addUser/addUserForm";
import {render, screen} from "@testing-library/react";
import {EmailChecker, PasswordChecker} from "@/__tests__/util/formCheckers";
import {permissionsList} from "@/components/permissions";

test("Add user form is rendered", () => {
    render(<AddUserForm/>)
    EmailChecker("Email Address")
    PasswordChecker()


    permissionsList.forEach(permission => {
        expect(screen.getByRole("checkbox", {name: permission.label})).toBeDefined();
    })
    expect(screen.getByRole("button", {name: "Create User"})).toBeDefined();
})