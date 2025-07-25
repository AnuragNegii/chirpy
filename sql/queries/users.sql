-- name: CreateUser :one
INSERT INTO users(id, created_at, updated_at, email, hashed_password)
VALUES(
    gen_random_UUID(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: ChangeUserNameAndPassword :exec
UPDATE users
    SET email = $1,
    hashed_password = $2,
    updated_at = NOW()
    WHERE id = $3;


-- name: GetUserFromUserID :one
SELECT * FROM users WHERE users.id = $1;
