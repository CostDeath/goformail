import {test} from "vitest";
import {render} from "@testing-library/react";
import LoginForm from "@/components/loginSignup/loginForm";
import {EmailChecker, PasswordChecker} from "@/__tests__/util/formCheckers";

test("Login form is rendered", () => {
    render(<LoginForm />);
    EmailChecker();
    PasswordChecker();
})