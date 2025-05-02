
export interface Email {
    id: number,
    rcpt: string[],
    sender: string,
    content: string,
    received_at: string,
    next_retry: string,
    exhausted: number,
    sent: boolean,
    list: string,
    approved: boolean
}