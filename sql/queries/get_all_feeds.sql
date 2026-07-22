-- name: GetFeeds :many
SELECT 
    f.name,
    f.url,
    u.name
FROM 
    feeds f
JOIN 
    users u ON f.user_id = u.id;
