DROP SEQUENCE IF EXISTS user_id_seq;
DROP TABLE IF EXISTS users;

CREATE SEQUENCE user_id_seq;
CREATE TABLE users
(
    id       integer NOT NULL DEFAULT nextval('user_id_seq'),
    login    varchar,
    name     varchar,
    password varchar,
    PRIMARY KEY (id)
);
ALTER SEQUENCE user_id_seq OWNED BY users.id;

DROP SEQUENCE IF EXISTS book_id_seq;
DROP TABLE IF EXISTS books;

CREATE SEQUENCE book_id_seq;
CREATE TABLE books
(
    id          integer NOT NULL DEFAULT nextval('book_id_seq'),
    name        varchar NOT NULL,
    author      varchar NOT NULL,
    description varchar NOT NULL DEFAULT '',
    holderId    integer NOT NULL DEFAULT 0,
    PRIMARY KEY (id)
);
ALTER SEQUENCE book_id_seq OWNED BY books.id;
