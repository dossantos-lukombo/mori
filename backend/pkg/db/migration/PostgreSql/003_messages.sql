-- +migrate Up
CREATE TABLE IF NOT EXISTS messages (
    message_id  varchar(100) NOT NULL,
    receiver_id varchar(100) NOT NULL, -- le groupe destinataire
    is_read     int4 DEFAULT 0,
    CONSTRAINT group_messages_pkey PRIMARY KEY (message_id),
    CONSTRAINT fk_group_messages_receiver FOREIGN KEY (receiver_id)
        REFERENCES groups(group_id) ON DELETE CASCADE
);

-- +migrate Down
DROP TABLE IF EXISTS messages;