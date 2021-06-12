DROP TABLE IF EXISTS forums CASCADE;
DROP TABLE IF EXISTS users CASCADE;

CREATE TABLE IF NOT EXISTS users
(
    nickname TEXT PRIMARY KEY,
    fullname TEXT NOT NULL,
    about    TEXT,
    email    TEXT UNIQUE
);

CREATE TABLE IF NOT EXISTS forums
(
    title        TEXT NOT NULL,
    userNickname TEXT NOT NULL,
    slug         TEXT UNIQUE,
    posts        INTEGER,
    threads      INTEGER
);