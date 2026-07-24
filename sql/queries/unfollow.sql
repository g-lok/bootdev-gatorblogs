-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows ff
USING feeds f
WHERE ff.feed_id = f.id
  AND ff.user_id = $1
  AND f.url = $2;


