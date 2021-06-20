DROP TABLE IF EXISTS forums CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS threads CASCADE;
DROP TABLE IF EXISTS posts CASCADE;
DROP TABLE IF EXISTS votes CASCADE;

CREATE UNLOGGED TABLE IF NOT EXISTS users
(
    id       SERIAL PRIMARY KEY,
    nickname CITEXT UNIQUE NOT NULL,
    fullname TEXT          NOT NULL,
    about    TEXT,
    email    CITEXT UNIQUE
);

CREATE UNLOGGED TABLE IF NOT EXISTS forums
(
    id           SERIAL PRIMARY KEY,
    title        TEXT                                                 NOT NULL,
    userNickname CITEXT REFERENCES users (nickname) ON DELETE CASCADE NOT NULL,
    slug         CITEXT UNIQUE,
    posts        INTEGER,
    threads      INTEGER
);

CREATE UNLOGGED TABLE IF NOT EXISTS threads
(
    id      SERIAL PRIMARY KEY,
    title   TEXT,
    author  CITEXT REFERENCES users (nickname) ON DELETE CASCADE NOT NULL,
    forum   CITEXT REFERENCES forums (slug) ON DELETE CASCADE,
    message TEXT, -- описание ветки
    votes   INTEGER DEFAULT 0                                    NOT NULL,
    slug    CITEXT                                               NOT NULL,
    created TIMESTAMP with time zone
);

CREATE UNLOGGED TABLE IF NOT EXISTS posts
(
    id       SERIAL PRIMARY KEY,
    parent   INTEGER DEFAULT NULL,
    forum    CITEXT REFERENCES forums (slug) ON DELETE CASCADE    NOT NULL,
    author   CITEXT REFERENCES users (nickname) ON DELETE CASCADE NOT NULL,
    thread   INTEGER REFERENCES threads (id) ON DELETE CASCADE    NOT NULL,
    created  TIMESTAMP with time zone,
    message  TEXT,
    isEdited BOOLEAN DEFAULT FALSE
);

CREATE UNLOGGED TABLE votes
(
    id     SERIAL PRIMARY KEY,
    author CITEXT REFERENCES users (nickname) ON DELETE CASCADE NOT NULL,
    thread INTEGER REFERENCES threads (id) ON DELETE CASCADE    NOT NULL,
    voice  INTEGER                                              NOT NULL,
    UNIQUE (author, thread)
);

select *
from users;
select *
from forums;
select *
from threads;
select *
from posts;

SELECT COUNT(*)
from pg_stat_activity;
SELECT *
FROM pg_stat_activity;
