-- name: FollowFeed :one
INSERT INTO feed_follows (id, feed_id, user_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetFeedFollowsByUserApiKey :many
SELECT * FROM feed_follows
WHERE user_id = (SELECT id FROM users WHERE apikey = $1);

-- name: DeleteFeedFollowByFeedFollowID :exec
DELETE FROM feed_follows
WHERE id = $1;