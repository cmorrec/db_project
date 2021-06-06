DROP TABLE IF EXISTS users CASCADE;

CREATE TABLE IF NOT EXISTS Users
(
    nickname TEXT PRIMARY KEY,
    fullname TEXT NOT NULL,
    about    TEXT,
    email    TEXT UNIQUE
);