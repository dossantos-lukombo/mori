-- +migrate Up
CREATE TABLE IF NOT EXISTS group_users (
    group_id varchar(100) NOT NULL,
    user_id  varchar(100) NOT NULL,
    CONSTRAINT pk_group_users PRIMARY KEY (group_id, user_id),
    CONSTRAINT fk_group_users_group FOREIGN KEY (group_id)
        REFERENCES groups(group_id) ON DELETE CASCADE,
    CONSTRAINT fk_group_users_user FOREIGN KEY (user_id)
        REFERENCES users(user_id) ON DELETE CASCADE
);

-- +migrate Down
DROP TABLE IF EXISTS group_users;