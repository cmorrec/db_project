CREATE EXTENSION citext;

DROP TABLE forums_users CASCADE;
DROP TABLE votes CASCADE;
DROP TABLE posts CASCADE;
DROP TABLE threads CASCADE;
DROP TABLE users CASCADE;
DROP TABLE forums CASCADE;

CREATE UNLOGGED TABLE users (
    nickname CITEXT UNIQUE NOT NULL COLLATE "POSIX",
    fullname TEXT,
    about TEXT,
    email CITEXT UNIQUE
);

CREATE INDEX IF NOT EXISTS user_nickname ON users using hash (nickname);
CREATE INDEX IF NOT EXISTS user_email ON users using hash (email);
CREATE INDEX IF NOT EXISTS user_all ON users (nickname, fullname, about, email);


CREATE UNLOGGED TABLE forums (
    id SERIAL PRIMARY KEY,
    user_create CITEXT REFERENCES users(nickname) ON DELETE CASCADE NOT NULL,
    title TEXT,
    slug CITEXT UNIQUE NOT NULL,
    threads INTEGER DEFAULT 0 NOT NULL,
    posts INTEGER DEFAULT 0 NOT NULL
);

CREATE INDEX IF NOT EXISTS forum_slug ON forums using hash (slug);


CREATE UNLOGGED TABLE threads (
    id SERIAL PRIMARY KEY,
    title TEXT,
    user_create CITEXT REFERENCES users(nickname) ON DELETE CASCADE NOT NULL,
    forum CITEXT REFERENCES forums(slug) ON DELETE CASCADE ,
    message TEXT, -- описание ветки
    votes INTEGER DEFAULT 0 NOT NULL,
    slug CITEXT NOT NULL,
    created TIMESTAMP with time zone
);

CREATE INDEX IF NOT EXISTS thr_slug ON threads using hash (slug);
CREATE INDEX IF NOT EXISTS thr_forum_created on threads (forum, created);


CREATE UNLOGGED TABLE posts (
    id SERIAL PRIMARY KEY,
    title TEXT,
    root_id INTEGER NOT NULL,
    parent INTEGER REFERENCES posts(id) DEFAULT NULL,
    forum CITEXT REFERENCES forums(slug) ON DELETE CASCADE NOT NULL,
    user_create CITEXT REFERENCES users(nickname) ON DELETE CASCADE NOT NULL,
    thread INTEGER REFERENCES threads(id) ON DELETE CASCADE NOT NULL,
    created TIMESTAMP with time zone,
    message TEXT,
    is_edited BOOLEAN DEFAULT FALSE,
    path INTEGER[]
);

create index idx_posts_thread on posts (thread);
create index idx_posts_path on posts using gin (path);
create index idx_posts_root_id on posts (root_id);
create index idx_posts_forum on posts (forum);


CREATE UNLOGGED TABLE votes (
    id SERIAL PRIMARY KEY,
    user_create CITEXT REFERENCES users(nickname) ON DELETE CASCADE NOT NULL,
    thread INTEGER REFERENCES threads(id) ON DELETE CASCADE NOT NULL,
    voice INTEGER NOT NULL,
    UNIQUE (user_create, thread)
);

CREATE UNLOGGED TABLE forums_users (
    user_nickname CITEXT REFERENCES users(nickname) ON DELETE CASCADE NOT NULL COLLATE "POSIX",
    user_fullname TEXT,
    user_about TEXT,
    user_email CITEXT,
    forum CITEXT REFERENCES forums(slug) ON DELETE CASCADE NOT NULL,
    UNIQUE (user_nickname, forum)
);

CREATE INDEX IF NOT EXISTS forums_users_forum_nickname on forums_users (forum, user_nickname);


CREATE OR REPLACE FUNCTION add_path() RETURNS TRIGGER AS
$add_path$
declare
    parents INTEGER[];
begin
    if (new.parent is null) then
        new.path := new.path || new.id;
        new.root_id := new.path[1];
    else
        select path from posts where id = new.parent and thread = new.thread
        into parents;

        if (coalesce(array_length(parents, 1), 0) = 0) then
            raise exception 'parent post not exists' USING ERRCODE = '12345';
        end if;

        new.path := new.path || parents || new.id;
        new.root_id := new.path[1];
    end if;
    return new;
end;
$add_path$ LANGUAGE plpgsql;

create trigger add_path
    before insert on posts for each row
execute procedure add_path();

-- функция и триггер при создании поста, на увеличение кол-ва постов в forums
CREATE OR REPLACE FUNCTION insert_post() RETURNS TRIGGER AS
$insert_post$
BEGIN
    UPDATE forums SET posts=posts + 1 WHERE forums.slug = NEW.forum;
    RETURN NEW;
END
$insert_post$ LANGUAGE plpgsql;

CREATE TRIGGER insert_post
AFTER INSERT ON posts
    FOR EACH ROW EXECUTE PROCEDURE insert_post();


-- функция и триггер при создании ветки, на увеличение кол-ва веток в forums
CREATE OR REPLACE FUNCTION insert_thread() RETURNS TRIGGER AS
$insert_thread$
BEGIN
    UPDATE forums SET threads=threads + 1 WHERE forums.slug = NEW.forum;
    RETURN NEW;
END
$insert_thread$ LANGUAGE plpgsql;

CREATE TRIGGER insert_thread
AFTER INSERT ON threads
    FOR EACH ROW EXECUTE PROCEDURE insert_thread();


-- функция и триггер при создании ветки и поста, на добавления пользователя в список форума
CREATE OR REPLACE FUNCTION new_forum_user_added() RETURNS TRIGGER AS
$new_forum_user_added$
BEGIN
    DECLARE
        nickAuthor citext;
        fullnameAuthor text;
        emailAuthor citext;
        aboutAuthor text;
    BEGIN
        SELECT nickname, fullname, about, email
        FROM users WHERE nickname = NEW.user_create
        INTO nickAuthor, fullnameAuthor, aboutAuthor, emailAuthor;

        INSERT INTO forums_users(user_fullname, user_about, user_email, user_nickname, forum)
        VALUES (fullnameAuthor, aboutAuthor, emailAuthor, nickAuthor, new.forum)
        ON CONFLICT DO NOTHING;

        RETURN NULL;
    END;
END;
$new_forum_user_added$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS new_forum_user_added ON posts;
CREATE TRIGGER new_forum_user_added
    AFTER INSERT ON posts
    FOR EACH ROW EXECUTE PROCEDURE new_forum_user_added();

DROP TRIGGER IF EXISTS new_forum_user_added ON threads;
CREATE TRIGGER new_forum_user_added
    AFTER INSERT ON threads
    FOR EACH ROW EXECUTE PROCEDURE new_forum_user_added();

-- функция и триггер при создании голоса, на увеличение кол-ва голосов в threads
CREATE OR REPLACE FUNCTION insert_voice() RETURNS TRIGGER AS
$insert_voice$
BEGIN
    UPDATE threads SET votes=(votes + NEW.voice) WHERE id = NEW.thread;
    RETURN NULL;
END
$insert_voice$ LANGUAGE plpgsql;

CREATE TRIGGER insert_vote
AFTER INSERT ON votes
    FOR EACH ROW EXECUTE PROCEDURE insert_voice();


-- функция и триггер при обновлении голоса, на изменение кол-ва голосов в threads
CREATE OR REPLACE FUNCTION update_voice() RETURNS TRIGGER AS
$update_voice$
BEGIN
    UPDATE threads SET votes= votes - OLD.voice + NEW.voice  WHERE id = NEW.thread;

    RETURN NULL;
END
$update_voice$ LANGUAGE plpgsql;

CREATE TRIGGER update_voice
AFTER UPDATE ON votes
    FOR EACH ROW EXECUTE PROCEDURE update_voice();


