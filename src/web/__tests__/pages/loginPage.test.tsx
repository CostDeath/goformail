import {test} from "vitest";
import {render} from "@testing-library/react";
import Page from "@/app/page";
import {EmailChecker, PasswordChecker} from "@/__tests__/util/formCheckers";

test("Login page is rendered", () => {
    render(<Page />);
    EmailChecker();

    PasswordChecker();
})
