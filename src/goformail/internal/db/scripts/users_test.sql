CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    hash TEXT NOT NULL,
    permissions TEXT[]
);

INSERT INTO users (email, permissions, hash)
VALUES ('get@test-0.tld', ARRAY['ADMIN'], 'hash');    -- ID: 1
INSERT INTO users (email, permissions, hash)
VALUES ('update@test-0.tld', ARRAY['ADMIN'], 'hash');  -- ID: 2
INSERT INTO users (email, permissions, hash)
VALUES ('update@test-1.tld', ARRAY['ADMIN'], 'hash');  -- ID: 3
INSERT INTO users (email, permissions, hash)
VALUES ('delete@test-0.tld', ARRAY['ADMIN'], 'hash');  -- ID: 4
