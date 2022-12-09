-- name: InsertRequisition :one
INSERT INTO requisitions (requisition, data, reasons, username, tenant, type, status, pubkey)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: ListRequisitions :many
SELECT *
FROM requisitions
WHERE tenant = $1
  AND username = $2
  AND (
    case
        when sqlc.arg('status') = '' then status = ANY (enum_range(NULL::requisition_status, null))
        else status = ANY (enum_range(sqlc.arg('status')::requisition_status, null)) end
    );

-- name: UpdateRequisition :exec
UPDATE requisitions
SET reasons    = $1,
    status     = $2,
    data       = $3,
    pubkey     = (
        case
            when sqlc.arg('pubkey') = '' OR sqlc.arg('pubkey') is null then pubkey
            else sqlc.arg('pubkey') end
        ),
    updated_at = NOW()
WHERE requisition = $4;

-- name: FailRequisition :exec
UPDATE requisitions
SET reasons    = (
    case
        when sqlc.arg('reasons') = '' then reasons
        else sqlc.arg('reasons') end
    ),
    status     = 'failure',
    updated_at = NOW()
WHERE requisition = $1;

-- name: GetRequisition :one
SELECT *
FROM requisitions
WHERE requisition = $1;

-- name: GetRequisitionById :one
SELECT *
FROM requisitions
WHERE id = $1;

-- name: RetryRequisition :exec
UPDATE requisitions
SET "retryTimes" = "retryTimes" + 1
WHERE requisition = $1;

-- name: GetRetryTimes :one
SELECT "retryTimes"
FROM requisitions
WHERE requisition = $1;
