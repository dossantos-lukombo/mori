-- +migrate Up
CREATE TABLE IF NOT EXISTS notifications (
    notif_id varchar(100) NOT NULL,
    user_id  varchar(100) NOT NULL,
    type   varchar(50)  NOT NULL,
    content varchar(255) NOT NULL,
    sender   varchar(100) NOT NULL,
    CONSTRAINT notifications_pkey PRIMARY KEY (notif_id),
    CONSTRAINT fk_notifications_user   FOREIGN KEY (user_id)
        REFERENCES public.users(user_id) ON DELETE CASCADE,
    CONSTRAINT fk_notifications_sender FOREIGN KEY (sender)
        REFERENCES public.users(user_id) ON DELETE CASCADE
);


-- +migrate Down
DROP TABLE IF EXISTS notifications;