-- +goose Up
ALTER TABLE feeds
ADD CONSTRAINT url_unique UNIQUE (url);

-- +goose Down
ALTER TABLE feeds
DROP CONSTRAINT url_unique;