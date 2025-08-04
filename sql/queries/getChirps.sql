-- name: GetChirps :many
SELECT * FROM chirps
ORDER BY created_at ASC;

-- name: GetChirpByID :one
SELECT * FROM chirps
WHERE id = $1;

-- name: GetChirpByAuthorId :many
SELECT * FROM chirps
WHERE user_id=$1;
