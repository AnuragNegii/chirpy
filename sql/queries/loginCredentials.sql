-- name: CheckLogin :one
SELECT * FROM users
where email = $1;
