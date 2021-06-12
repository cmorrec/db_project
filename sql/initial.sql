DROP TABLE IF EXISTS forums CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS threads CASCADE;

CREATE TABLE IF NOT EXISTS users
(
    id       SERIAL PRIMARY KEY,
    nickname CITEXT UNIQUE NOT NULL,
    fullname TEXT          NOT NULL,
    about    TEXT,
    email    CITEXT UNIQUE
);

CREATE TABLE IF NOT EXISTS forums
(
    id           SERIAL PRIMARY KEY,
    title        TEXT                                               NOT NULL,
    userNickname CITEXT REFERENCES users (nickname) ON DELETE CASCADE NOT NULL,
    slug         CITEXT UNIQUE,
    posts        INTEGER,
    threads      INTEGER
);

CREATE TABLE IF NOT EXISTS threads
(
    id      SERIAL PRIMARY KEY,
    title   TEXT,
    author  TEXT REFERENCES users (nickname) ON DELETE CASCADE NOT NULL,
    forum   TEXT REFERENCES forums (slug) ON DELETE CASCADE,
    message TEXT, -- описание ветки
    votes   INTEGER DEFAULT 0                                  NOT NULL,
    slug    TEXT                                               NOT NULL,
    created TIMESTAMP with time zone
);