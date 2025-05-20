-- +migrate Up
CREATE TABLE IF NOT EXISTS conversations (
    id               serial4       NOT NULL,
    user_id          varchar(100)  NOT NULL,
    conversation_id  varchar(100)  NOT NULL,
    convo            _jsonb        NOT NULL,
    new_conversation bool DEFAULT false NOT NULL,
    CONSTRAINT conversations_pkey PRIMARY KEY (id),
    CONSTRAINT fk_conversations_user FOREIGN KEY (user_id)
        REFERENCES users(user_id) ON DELETE CASCADE
);

-- +migrate Down
DROP TABLE IF EXISTS conversations;
