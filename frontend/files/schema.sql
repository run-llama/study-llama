-- Files table
CREATE TABLE IF NOT EXISTS files (
    id SERIAL PRIMARY KEY,
    username TEXT NOT NULL,
    file_name TEXT NOT NULL,
    file_category TEXT DEFAULT NULL
);