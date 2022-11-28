-- name: CreateShare :one
INSERT INTO shares ( pubkey, data, enable, notification, address)
VALUES ($1, $2, true, 'enable', LOWER(sqlc.arg('address'))) RETURNING *;

-- name: GetShare :one
SELECT *
FROM shares
WHERE id = $1 LIMIT 1;

-- name: GetShareByAddress :one
SELECT *
FROM shares
WHERE LOWER(address) = LOWER(sqlc.arg('address')) LIMIT 1;

-- name: GetShareByPubkey :one
SELECT *
FROM shares
WHERE LOWER(pubkey) = LOWER(sqlc.arg('pubkey')) LIMIT 1;

-- name: ListShare :many
SELECT *
FROM shares
ORDER BY id LIMIT $1
OFFSET $2;
