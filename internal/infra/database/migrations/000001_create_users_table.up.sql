CREATE TABLE IF NOT EXISTS public.users (
    id              serial PRIMARY KEY,
    email           varchar(50) NOT NULL,
    password        varchar(100) NOT NULL,
    "name"          varchar(50) NOT NULL,
    phone_number    varchar(30) NOT NULL,
    avatar          varchar NULL,
    activated       boolean NOT NULL,
    created_date    timestamp NOT NULL,
    updated_date    timestamp NOT NULL,
    deleted_date    timestamp NULL
);
