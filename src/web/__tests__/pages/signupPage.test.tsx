import {expect, test} from "vitest";
import {render, screen} from "@testing-library/react";
import Page from "@/app/signup/page";
import {EmailChecker, PasswordChecker, StudentIDChecker} from "@/__tests__/util/formCheckers";


test("Whole page is rendered", () => {
    render(<Page />);

    EmailChecker();

    StudentIDChecker();

    PasswordChecker();


    expect(screen.getByTestId("to-login")).toBeDefined();
})