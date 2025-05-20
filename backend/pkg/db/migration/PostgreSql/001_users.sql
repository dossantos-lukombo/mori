-- +migrate Up
CREATE TABLE IF NOT EXISTS users (
    user_id             varchar(100) NOT NULL,
    created_at          timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    email               varchar(50)  NOT NULL,
    first_name          varchar(50)  NOT NULL,
    last_name           varchar(50)  NOT NULL,
    nickname            varchar(50),
    birthday            date         NOT NULL,
    image               varchar(255),
    about               text,
    status              varchar(50) DEFAULT 'PUBLIC' NOT NULL,
    password         varchar(100) NOT NULL,
    verification_token  varchar(100),
    verified            bool DEFAULT false NOT NULL,
    reset_token         varchar(100),
    reset_token_expires timestamp,
    CONSTRAINT users_pkey PRIMARY KEY (user_id)
);

-- +migrate Down
DROP TABLE IF EXISTS users;
