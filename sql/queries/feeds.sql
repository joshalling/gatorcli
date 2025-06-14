-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds inner join users on feeds.user_id = users.id;

-- name: GetFeedByUrl :one
SELECT * FROM feeds
WHERE url = $1 LIMIT 1;