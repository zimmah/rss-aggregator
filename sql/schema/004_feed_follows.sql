-- +goose Up
CREATE TABLE feed_follows (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    user_id UUID NOT NULL,
    CONSTRAINT fk_user
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    feed_id UUID NOT NULL,
    CONSTRAINT fk_feed
    FOREIGN KEY (feed_id) REFERENCES feeds(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feed_follows;