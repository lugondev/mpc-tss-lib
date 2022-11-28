-- name: MapEmailContract :one
WITH inserted_data AS (
    INSERT
        INTO emails_contract
            (email_id, contract_id)
            VALUES ($1, $2) RETURNING *)
SELECT *
FROM inserted_data
         INNER JOIN emails ON inserted_data.email_id = emails.id
         INNER JOIN contracts ON inserted_data.contract_id = contracts.id
WHERE contracts.id = inserted_data.contract_id
  AND emails.id = inserted_data.email_id;

-- name: ListEmailsSubscription :many
SELECT emails.*
FROM emails_contract
         JOIN emails ON emails_contract.email_id = emails.id
WHERE contract_id = sqlc.arg('contractId');

-- name: ListEmailsSubscriptionByAddress :many
SELECT emails.*
FROM contracts
         JOIN emails_contract ON emails_contract.contract_id = contracts.id
         JOIN emails ON emails_contract.email_id = emails.id
WHERE contracts.address = sqlc.arg('contractAddress');
