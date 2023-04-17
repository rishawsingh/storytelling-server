ALTER TABLE IF EXISTS users ALTER COLUMN password DROP NOT NULL;

CREATE TABLE IF NOT EXISTS sessions
(
    id          TEXT PRIMARY KEY,
    user_id     INTEGER REFERENCES users (id),
    archived_at TIMESTAMP WITH TIME ZONE
);
