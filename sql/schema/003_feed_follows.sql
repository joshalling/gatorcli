-- +goose Up
CREATE TABLE feed_follows (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW() NOT NULL,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id),
    feed_id UUID NOT NULL REFERENCES feeds (id) ON DELETE CASCADE,
    FOREIGN KEY (feed_id) REFERENCES feeds (id),
    UNIQUE (user_id, feed_id)
);

-- +goose Down
DROP TABLE feed_follows;