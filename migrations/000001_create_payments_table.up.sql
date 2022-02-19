create table payments
(
    id           bigserial
        constraint payments_pk
            primary key,
    name         text                       not null,
    date         timestamp(0) default now() not null,
    payment_type text                       not null,
    comment      text                       not null,
    category_id  bigint                     not null
        constraint payments_categories_id_fk
            references categories,
    price        integer                    not null
);
