CREATE TYPE "role_type" AS ENUM (
    'admin',
    'user'
    );

alter table if exists public.users
    add role role_type NOT NULL