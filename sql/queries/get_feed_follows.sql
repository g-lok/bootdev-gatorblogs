-- name: GetFeedFollowsForUser :many
SELECT 
  ff.id,
  ff.user_id,
  ff.feed_id,
  ff.created_at,
  ff.updated_at,
  u.name AS user_name,
  f.name AS feed_name
FROM feed_follows ff
JOIN users u ON u.id = ff.user_id
JOIN feeds f ON f.id = ff.feed_id
WHERE ff.user_id = $1;
