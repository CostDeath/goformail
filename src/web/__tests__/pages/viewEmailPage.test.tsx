import {expect, test} from "vitest";
import {render} from "@testing-library/react";
import Page from "@/app/(dashboards)/mailingLists/list/email/page"
import {EmailViewChecker} from "@/__tests__/util/emailViewChecker";

test("Email view page is rendered", () => {
    render(<Page />);

    EmailViewChecker()
})