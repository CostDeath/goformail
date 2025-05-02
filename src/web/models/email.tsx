
export interface Email {
    id: number,
    rcpt: string[],
    sender: string,
    content: string,
    received_at: Date,
    next_retry: Date,
    exhausted: number,
    sent: boolean,
    list: string,
    approved: boolean
}