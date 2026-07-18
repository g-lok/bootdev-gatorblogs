-- queries.sql
-- name: UserExists :one
SELECT EXISTS(
    SELECT 1 FROM users WHERE name = @name
);
