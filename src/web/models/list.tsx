
export interface List {
    name: string,
    recipients: string[]
}

export interface MailingLists {
    id: number,
    list: List
}
