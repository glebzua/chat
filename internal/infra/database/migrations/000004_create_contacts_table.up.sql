CREATE TABLE IF NOT EXISTS public.contacts (
    id              serial PRIMARY KEY,
    userId          int NOT NULL,
    contactId       int NOT NULL,
    activated       boolean NOT NULL,
    chatId char(60),
    created_date    timestamp NOT NULL,
    updated_date    timestamp NOT NULL,
    deleted_date    timestamp NULL
);
