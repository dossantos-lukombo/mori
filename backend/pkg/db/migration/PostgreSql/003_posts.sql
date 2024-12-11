-- +migrate Up
CREATE TABLE IF NOT EXISTS posts (
    post_id VARCHAR(100) NOT NULL PRIMARY KEY,
    group_id VARCHAR(100),
    created_by VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    content TEXT,
    image VARCHAR(50),
    visibility VARCHAR(50) DEFAULT 'PUBLIC'
);

-- +migrate Down
DROP TABLE IF EXISTS posts;