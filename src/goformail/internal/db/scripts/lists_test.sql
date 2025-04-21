CREATE TABLE lists (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    recipients TEXT[]
);

INSERT INTO lists (name, recipients) VALUES ('get-test-0', ARRAY['example@domain.tld']);    -- ID: 1
INSERT INTO lists (name, recipients) VALUES ('patch-test-0', ARRAY['example@domain.tld']);  -- ID: 2
INSERT INTO lists (name, recipients) VALUES ('patch-test-1', ARRAY['example@domain.tld']);  -- ID: 3
INSERT INTO lists (name, recipients) VALUES ('delete-test-0', ARRAY['example@domain.tld']);  -- ID: 4
