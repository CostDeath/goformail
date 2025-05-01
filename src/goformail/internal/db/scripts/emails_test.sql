CREATE TABLE emails (
    id SERIAL PRIMARY KEY,
    rcpt TEXT[] NOT NULL,
    sender TEXT NOT NULL,
    content TEXT NOT NULL,
    received_at TIMESTAMPTZ NOT NULL,
    next_retry TIMESTAMPTZ,
    exhausted INT DEFAULT 3,
    sent BOOL DEFAULT false,
    list INT,
    approved BOOL DEFAULT false
);

CREATE TABLE lists (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    recipients TEXT[],
    mods INT[],
    approved_senders TEXT[],
    locked BOOL DEFAULT false
);

INSERT INTO emails (rcpt, sender, content, received_at, next_retry, exhausted, sent, list, approved)
VALUES (ARRAY['test-1@test.tld'], 'sender@test-0.tld', 'content', '2025-04-30 15:00:00', '2025-04-30 15:00:00', 2, false, 1, true);    -- ID: 1
INSERT INTO emails (rcpt, sender, content, received_at, next_retry, exhausted, sent, list, approved)
VALUES (ARRAY['test-2@test.tld'], 'sender@test-0.tld', 'content', '2025-04-30 15:00:00', '2025-04-30 15:00:00', 2, false, 1, false);    -- ID: 2
INSERT INTO emails (rcpt, sender, content, received_at, next_retry, exhausted, sent, list, approved)
VALUES (ARRAY['test-3@test.tld'], 'sender@test-0.tld', 'content', '2025-04-30 15:00:00', '2025-04-30 15:00:00', 0, false, 1, true);    -- ID: 3
INSERT INTO emails (rcpt, sender, content, received_at, next_retry, exhausted, sent, list, approved)
VALUES (ARRAY['test-4@test.tld'], 'sender@test-0.tld', 'content', '2025-04-30 15:00:00', '3025-04-30 15:00:00', 3, false, 2, true);    -- ID: 4
INSERT INTO emails (rcpt, sender, content, received_at, next_retry, exhausted, sent, list, approved)
VALUES (ARRAY['test-5@test.tld'], 'sender@test-0.tld', 'content', '2025-04-30 15:00:00', '3025-04-30 15:00:00', 3, false, 2, true);    -- ID: 5

INSERT INTO lists (name, recipients, mods, approved_senders)
VALUES ('test-0', ARRAY['example@domain.tld'], ARRAY[1], ARRAY['sender@test-0.tld']);    -- ID: 1
