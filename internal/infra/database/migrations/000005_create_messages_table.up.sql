CREATE TABLE IF NOT EXISTS public.messages (
    id              serial PRIMARY KEY,
    chatId          char(60) NOT NULL,
    senderId        int NOT NULL,
    recipientId     int NOT NULL,
    message         varchar(350) NOT NULL DEFAULT '',
    fileloc VARCHAR (250) NOT NULL DEFAULT '',
    sended          boolean NOT NULL,
    received         boolean NOT NULL,
    created_date    timestamp NOT NULL,
    updated_date    timestamp NOT NULL,
    deleted_date    timestamp NULL
);
