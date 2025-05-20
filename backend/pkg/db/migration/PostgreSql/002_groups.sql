-- +migrate Up
CREATE TABLE IF NOT EXISTS groups (
    group_id      varchar(100) NOT NULL,
    administrator varchar(100) NOT NULL,
    name        varchar(50)  NOT NULL,
    description   varchar(255),
    CONSTRAINT groups_pkey PRIMARY KEY (group_id),
    -- L’administrateur est un user
    CONSTRAINT fk_groups_admin FOREIGN KEY (administrator)
        REFERENCES users(user_id) ON DELETE CASCADE
);

-- +migrate Down
DROP TABLE IF EXISTS groups;