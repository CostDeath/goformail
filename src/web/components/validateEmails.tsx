
export default function validateEmail(email: string) {
    const regex = new RegExp("^([A-z0-9+._/&!][-A-z0-9+._/&!]*)@(([a-z0-9][-a-z0-9]*\.)([-a-z0-9]+\.)*[a-z]{2,})$")
    return regex.test(email)
}