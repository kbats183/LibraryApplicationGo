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

insert into users(login, name, password)
values ('user', 'User', '40bd001563085fc35165329ea1ff5c5ecbdbbeef'); /* password is 123 */

insert into books(name, author, description)
values ('Война и мир', 'Лев Толстой',
        'Роман-эпопея Льва Николаевича Толстого, описывающий русское общество в эпоху войн против Наполеона в 1805-1812 годах.'),
       ('1984', 'Джордж Оруэлл', 'Роман-антиутопия Джорджа Оруэлла, изданный в 1949 году.'),
       ('Улисс', 'Джеймс Джойс', ''),
       ('Лолита', 'Владимир Набоков', ''),
       ('Шум и ярость', 'Уильям Фолкнер', '');
