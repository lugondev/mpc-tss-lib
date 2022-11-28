-- name: GetPrivateKey :one
SELECT *
FROM private_keys
WHERE LOWER(sqlc.arg('address')) = LOWER(address) LIMIT 1;

-- name: ListAddresses :many
SELECT address
FROM private_keys;

-- name: AddPrivateKey :one
INSERT INTO private_keys (private_key, address)
VALUES ($1, LOWER(sqlc.arg('address'))) RETURNING *;

-- name: GetChain :one
SELECT *
FROM chains
WHERE chain_id = $1 LIMIT 1;

-- name: AddChain :one
INSERT INTO chains (name, rpcs, chain_id)
VALUES ($1, $2, $3) RETURNING *;

-- name: ListChains :many
SELECT *
FROM chains;
