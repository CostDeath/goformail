CREATE TABLE lists (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    recipients TEXT[],
    mods INT[],
    approved_senders TEXT[],
    locked BOOL DEFAULT false
);

INSERT INTO lists (name, recipients, mods, approved_senders)
VALUES ('get-test-0', ARRAY['example@domain.tld'], ARRAY[1], ARRAY['example@domain.tld']);    -- ID: 1

INSERT INTO lists (name, recipients, mods, approved_senders)
VALUES ('patch-test-0', ARRAY['example@domain.tld'], ARRAY[1], ARRAY['example@domain.tld']);  -- ID: 2

INSERT INTO lists (name, recipients, mods, approved_senders)
VALUES ('patch-test-1', ARRAY['example@domain.tld'], ARRAY[1], ARRAY['example@domain.tld']);  -- ID: 3

INSERT INTO lists (name, recipients, mods, approved_senders)
VALUES ('patch-test-2', ARRAY['example@domain.tld'], ARRAY[1], ARRAY['example@domain.tld']);  -- ID: 4

INSERT INTO lists (name, recipients, mods, approved_senders, locked)
VALUES ('patch-test-3', ARRAY['example@domain.tld'], ARRAY[1], ARRAY['example@domain.tld'], true);  -- ID: 5

INSERT INTO lists (name, recipients, mods, approved_senders)
VALUES ('delete-test-0', ARRAY['example@domain.tld'], ARRAY[1], ARRAY['example@domain.tld']);  -- ID: 6

INSERT INTO lists (name, recipients, mods, approved_senders, locked)
VALUES ('patch-test-4', ARRAY['example@domain.tld'], ARRAY[1], ARRAY['example@domain.tld'], true);  -- ID: 7
