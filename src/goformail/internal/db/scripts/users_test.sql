CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    hash TEXT NOT NULL,
    permissions TEXT[]
);

CREATE TABLE lists (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    recipients TEXT[],
    mods INT[],
    approved_senders TEXT[],
    locked BOOL DEFAULT false
);

INSERT INTO users (email, permissions, hash)
VALUES ('get@test-0.tld', ARRAY['ADMIN'], 'hash');    -- ID: 1
INSERT INTO users (email, permissions, hash)
VALUES ('update@test-0.tld', ARRAY['ADMIN'], 'hash');  -- ID: 2
INSERT INTO users (email, permissions, hash)
VALUES ('update@test-1.tld', ARRAY['ADMIN'], 'hash');  -- ID: 3
INSERT INTO users (email, permissions, hash)
VALUES ('delete@test-0.tld', ARRAY['ADMIN'], 'hash');  -- ID: 4

INSERT INTO lists (name, recipients, mods, approved_senders)
VALUES ('get-test-0', ARRAY['example@domain.tld'], ARRAY[1], ARRAY['example@domain.tld']);    -- ID: 1
