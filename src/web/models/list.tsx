
export interface List {
    name: string,
    recipients: string[]
}

export interface MailingList {
    id: number,
    name: string,
    locked: boolean,
    recipients: string[]
    mods: number[]
    approvedSenders: string[]
}
