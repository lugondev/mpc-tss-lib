-- name: AddEmail :one
INSERT INTO emails (name, email)
VALUES ($1, lower(sqlc.arg(email)::text))
RETURNING *;

-- name: GetEmail :one
SELECT *
FROM emails
WHERE id = $1
LIMIT 1;

-- name: ListEmails :many
SELECT *
FROM emails
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: UpdateEmail :one
UPDATE emails
SET name = $2
WHERE id = $1
RETURNING *;

-- name: DeleteEmail :exec
DELETE
FROM emails
WHERE id = $1;
