-- +migrate Up
CREATE TABLE IF NOT EXISTS sessions (
    session_id      varchar(100) NOT NULL,
    user_id         varchar(100) NOT NULL,
    expiration_time timestamp    NOT NULL,
    CONSTRAINT sessions_pkey PRIMARY KEY (session_id),
    CONSTRAINT fk_sessions_user FOREIGN KEY (user_id)
    REFERENCES users(user_id) ON DELETE CASCADE
);
-- +migrate Down
DROP TABLE sessions;