import {expect, test} from "vitest";
import validateEmail from "@/components/validateEmails";


test("Given email is valid", () => {
    expect(validateEmail("some-Email.2@sub.example.com")).toBe(true)
})

test("Given email is invalid", () => {
    expect(validateEmail("-some-Email.2@@sub.example.com")).toBe(false)
})