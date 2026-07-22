-- name: CreateFeedFollow :one
WITH inserted_follow AS (
  INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
  VALUES ($1, $2, $3, $4, $5)
  RETURNING id, user_id, feed_id, created_at, updated_at
)
SELECT
  f.id,
  f.user_id,
  f.feed_id,
  f.created_at,
  f.updated_at,
  u.name AS user_name,
  fd.name AS feed_name
FROM inserted_follow f
JOIN users u ON u.id = f.user_id
JOIN feeds fd ON fd.id = f.feed_id;
