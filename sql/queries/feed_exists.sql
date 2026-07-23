-- queries.sql
-- name: FeedExists :one
SELECT EXISTS(
    SELECT 1 FROM feeds WHERE url = @url
);
