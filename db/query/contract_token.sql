-- name: AddTokenToContract :one
INSERT INTO tokens_contract (contract_id, name, symbol, token, decimals)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: AddTokenToContractByAddress :one
INSERT
INTO tokens_contract (contract_id, name, symbol, token, decimals)
VALUES ((SELECT id as contract_id FROM contracts WHERE LOWER(address) = LOWER(sqlc.arg('address'))), $1, $2, LOWER(sqlc.arg('token')), $3)
RETURNING *;

-- name: ListTokenInContract :many
SELECT *
FROM tokens_contract
WHERE contract_id = $1;

-- name: ListTokenInContractByAddress :many
SELECT *
FROM tokens_contract
WHERE contract_id = (SELECT id
                     FROM contracts
                     WHERE LOWER(address) = LOWER(sqlc.arg('address')));

-- name: RemoveTokenFromContract :exec
DELETE
FROM tokens_contract
WHERE contract_id = $1
  AND LOWER(token) = LOWER(sqlc.arg('token'));

-- name: RemoveTokenFromContractByAddress :exec
DELETE
FROM tokens_contract
WHERE contract_id = (SELECT id
                     FROM contracts
                     WHERE LOWER(address) = LOWER(sqlc.arg('address')))
  AND LOWER(token) = LOWER(sqlc.arg('token'));
