-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW() NOT NULL,
    name VARCHAR(255) UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE users;