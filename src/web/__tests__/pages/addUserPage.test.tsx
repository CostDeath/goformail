import {expect, test} from "vitest";
import {render, screen} from "@testing-library/react";
import Page from "@/app/(dashboards)/management/add/page";
import {EmailChecker, PasswordChecker} from "@/__tests__/util/formCheckers";
import {permissionsList} from "@/components/permissions";

test("add user page is rendered", () => {
    render(<Page />);
    EmailChecker();

    PasswordChecker();

    permissionsList.forEach(permission => {
        expect(screen.getByRole("checkbox", {name: permission.label})).toBeDefined()
    })

})