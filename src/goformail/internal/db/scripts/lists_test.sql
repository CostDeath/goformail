CREATE TABLE lists (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    recipients TEXT[]
);
INSERT INTO lists (name, recipients) VALUES ('get-test-0', ARRAY['example@domain.tld']);
