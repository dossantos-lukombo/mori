-- +migrate Up
CREATE TABLE IF NOT EXISTS followers (
    follower_id varchar(100) NOT NULL,
    user_id     varchar(100) NOT NULL,
    CONSTRAINT pk_followers PRIMARY KEY (follower_id, user_id),
    -- Les deux protagonistes pointent vers users
    CONSTRAINT fk_followers_follower FOREIGN KEY (follower_id)
        REFERENCES users(user_id) ON DELETE CASCADE,
    CONSTRAINT fk_followers_user FOREIGN KEY (user_id)
        REFERENCES users(user_id) ON DELETE CASCADE
);

-- +migrate Down
DROP TABLE IF EXISTS followers;