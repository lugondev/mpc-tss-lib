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
