import {test} from "vitest";
import {render} from "@testing-library/react";
import SignUpForm from "@/components/loginSignup/signUpForm";
import {EmailChecker, PasswordChecker, StudentIDChecker} from "@/__tests__/util/formCheckers";

test("Sign up form is rendered", () => {
    render(<SignUpForm />);
    EmailChecker();
    StudentIDChecker();
    PasswordChecker();
})