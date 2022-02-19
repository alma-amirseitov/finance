create table if not exists categories
(
    id   bigserial
    constraint categories_pk
    primary key,
    name text
);