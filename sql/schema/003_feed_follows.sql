-- +goose Up
CREATE TABLE feed_follows (
    id UUID PRIMARY KEY,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    user_id UUID NOT NULL,
    feed_id UUID NOT NULL,
    CONSTRAINT user_feed_pair UNIQUE (user_id, feed_id)
);

ALTER TABLE feed_follows
ADD FOREIGN KEY (user_id)
REFERENCES users(id)
ON DELETE CASCADE;

ALTER TABLE feed_follows
ADD FOREIGN KEY (feed_id)
REFERENCES feeds(id)
ON DELETE CASCADE;

-- +goose Down
DROP TABLE feed_follows;