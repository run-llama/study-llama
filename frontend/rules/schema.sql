-- Rules table
CREATE TABLE IF NOT EXISTS rules (
    id SERIAL PRIMARY KEY,
    username TEXT NOT NULL,
    rule_name TEXT NOT NULL,
    rule_type TEXT NOT NULL,
    rule_description TEXT NOT NULL
);