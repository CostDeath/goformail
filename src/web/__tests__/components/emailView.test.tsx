import {test} from "vitest";
import {render} from "@testing-library/react";
import EmailView from "@/components/emailView";
import {EmailViewChecker} from "@/__tests__/util/emailViewChecker";

test("Email view is rendered", () => {
    render(<EmailView id={"1"} />);

    EmailViewChecker()
})