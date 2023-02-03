CREATE TABLE IF NOT EXISTS images
(
    id           serial PRIMARY KEY,
    senderid      int          NOT NULL,
    chatid       char(60) NOT NULL,
    recipientid      int  NOT NULL,
    objid       int,
    obj_type     varchar(50),
    name         varchar(250) NOT NULL,
    created_date timestamp    NOT NULL,
    updated_date timestamp    NOT NULL,
    deleted_date timestamp    NULL,
    CONSTRAINT fk_users
    FOREIGN KEY (senderid)
    REFERENCES users (id)
    );
