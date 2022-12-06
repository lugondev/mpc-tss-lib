-- name: CreateShare :one
INSERT INTO shares (pubkey, data, address, party_id)
VALUES ($1, $2, LOWER(sqlc.arg('address')), $3)
RETURNING pubkey, enable, notification, address, party_id;

-- name: GetShare :one
SELECT pubkey, enable, notification, address, party_id
FROM shares
WHERE id = $1
LIMIT 1;

-- name: GetShareByAddress :one
SELECT pubkey, enable, notification, address, party_id
FROM shares
WHERE LOWER(address) = LOWER(sqlc.arg('address'))
LIMIT 1;

-- name: GetShareByID :one
SELECT pubkey, data
FROM shares
WHERE party_id = $1
LIMIT 1;

-- name: GetPartyIdByPubkey :one
SELECT pubkey, party_id, address
FROM shares
WHERE LOWER(pubkey) = LOWER(sqlc.arg('pubkey'))
LIMIT 1;

-- name: ListShare :many
SELECT pubkey, enable, notification, address, party_id
FROM shares
ORDER BY id;
