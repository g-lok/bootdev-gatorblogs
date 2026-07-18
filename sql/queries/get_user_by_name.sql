-- name: GetUserName :one
SELECT * FROM users WHERE name = $1;
