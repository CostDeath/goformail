CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    hash TEXT NOT NULL,
    salt TEXT NOT NULL,
    token TEXT,
    permissions TEXT[]
);

INSERT INTO users (email, permissions, hash, salt)
VALUES ('get@test-0.tld', ARRAY['ADMIN'], 'hash', 'salt');    -- ID: 1
INSERT INTO users (email, permissions, hash, salt)
VALUES ('update@test-0.tld', ARRAY['ADMIN'], 'hash', 'salt');  -- ID: 2
INSERT INTO users (email, permissions, hash, salt)
VALUES ('update@test-1.tld', ARRAY['ADMIN'], 'hash', 'salt');  -- ID: 3
INSERT INTO users (email, permissions, hash, salt)
VALUES ('delete@test-0.tld', ARRAY['ADMIN'], 'hash', 'salt');  -- ID: 4
