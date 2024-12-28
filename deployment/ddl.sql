-- Dump your initial database ddl in this file.
-- It will be setup when running docker compose up for the first time.
CREATE SCHEMA your_schema;
SET search_path TO your_schema;

create table users
(
    id        uuid                  not null
        constraint users_pk
            primary key,
    name      text,
    email     text,
    password  text
);

alter table users
    owner to your_postgres_user;

create table sessions
(
    id         uuid      not null,
    user_id    uuid      not null
        constraint sessions_user___fk
            references your_schema.users,
    expires_at timestamp not null
);

alter table sessions
    owner to your_postgres_user;

