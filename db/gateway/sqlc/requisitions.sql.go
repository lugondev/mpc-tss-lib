// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: requisitions.sql

package db_gateway

import (
	"context"
)

const failRequisition = `-- name: FailRequisition :exec
UPDATE requisitions
SET reasons    = (
    case
        when $2 = '' then reasons
        else $2 end
    ),
    status     = 'failure',
    updated_at = NOW()
WHERE requisition = $1
`

type FailRequisitionParams struct {
	Requisition string      `json:"requisition"`
	Reasons     interface{} `json:"reasons"`
}

func (q *Queries) FailRequisition(ctx context.Context, arg FailRequisitionParams) error {
	_, err := q.db.ExecContext(ctx, failRequisition, arg.Requisition, arg.Reasons)
	return err
}

const getRequisition = `-- name: GetRequisition :one
SELECT id, requisition, pubkey, data, reasons, username, tenant, "retryTimes", type, status, created_at, updated_at
FROM requisitions
WHERE requisition = $1
`

func (q *Queries) GetRequisition(ctx context.Context, requisition string) (Requisition, error) {
	row := q.db.QueryRowContext(ctx, getRequisition, requisition)
	var i Requisition
	err := row.Scan(
		&i.ID,
		&i.Requisition,
		&i.Pubkey,
		&i.Data,
		&i.Reasons,
		&i.Username,
		&i.Tenant,
		&i.RetryTimes,
		&i.Type,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getRequisitionById = `-- name: GetRequisitionById :one
SELECT id, requisition, pubkey, data, reasons, username, tenant, "retryTimes", type, status, created_at, updated_at
FROM requisitions
WHERE id = $1
`

func (q *Queries) GetRequisitionById(ctx context.Context, id int64) (Requisition, error) {
	row := q.db.QueryRowContext(ctx, getRequisitionById, id)
	var i Requisition
	err := row.Scan(
		&i.ID,
		&i.Requisition,
		&i.Pubkey,
		&i.Data,
		&i.Reasons,
		&i.Username,
		&i.Tenant,
		&i.RetryTimes,
		&i.Type,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getRetryTimes = `-- name: GetRetryTimes :one
SELECT "retryTimes"
FROM requisitions
WHERE requisition = $1
`

func (q *Queries) GetRetryTimes(ctx context.Context, requisition string) (int32, error) {
	row := q.db.QueryRowContext(ctx, getRetryTimes, requisition)
	var retryTimes int32
	err := row.Scan(&retryTimes)
	return retryTimes, err
}

const insertRequisition = `-- name: InsertRequisition :one
INSERT INTO requisitions (requisition, data, reasons, username, tenant, type, status, pubkey)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, requisition, pubkey, data, reasons, username, tenant, "retryTimes", type, status, created_at, updated_at
`

type InsertRequisitionParams struct {
	Requisition string            `json:"requisition"`
	Data        []byte            `json:"data"`
	Reasons     string            `json:"reasons"`
	Username    string            `json:"username"`
	Tenant      string            `json:"tenant"`
	Type        RequisitionType   `json:"type"`
	Status      RequisitionStatus `json:"status"`
	Pubkey      string            `json:"pubkey"`
}

func (q *Queries) InsertRequisition(ctx context.Context, arg InsertRequisitionParams) (Requisition, error) {
	row := q.db.QueryRowContext(ctx, insertRequisition,
		arg.Requisition,
		arg.Data,
		arg.Reasons,
		arg.Username,
		arg.Tenant,
		arg.Type,
		arg.Status,
		arg.Pubkey,
	)
	var i Requisition
	err := row.Scan(
		&i.ID,
		&i.Requisition,
		&i.Pubkey,
		&i.Data,
		&i.Reasons,
		&i.Username,
		&i.Tenant,
		&i.RetryTimes,
		&i.Type,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listRequisitions = `-- name: ListRequisitions :many
SELECT id, requisition, pubkey, data, reasons, username, tenant, "retryTimes", type, status, created_at, updated_at
FROM requisitions
WHERE tenant = $1
  AND username = $2
  AND (
    case
        when $3 = '' then status = ANY (enum_range(NULL::requisition_status, null))
        else status = ANY (enum_range($3::requisition_status, null)) end
    )
`

type ListRequisitionsParams struct {
	Tenant   string      `json:"tenant"`
	Username string      `json:"username"`
	Status   interface{} `json:"status"`
}

func (q *Queries) ListRequisitions(ctx context.Context, arg ListRequisitionsParams) ([]Requisition, error) {
	rows, err := q.db.QueryContext(ctx, listRequisitions, arg.Tenant, arg.Username, arg.Status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Requisition{}
	for rows.Next() {
		var i Requisition
		if err := rows.Scan(
			&i.ID,
			&i.Requisition,
			&i.Pubkey,
			&i.Data,
			&i.Reasons,
			&i.Username,
			&i.Tenant,
			&i.RetryTimes,
			&i.Type,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const retryRequisition = `-- name: RetryRequisition :exec
UPDATE requisitions
SET "retryTimes" = "retryTimes" + 1
WHERE requisition = $1
`

func (q *Queries) RetryRequisition(ctx context.Context, requisition string) error {
	_, err := q.db.ExecContext(ctx, retryRequisition, requisition)
	return err
}

const updateRequisition = `-- name: UpdateRequisition :exec
UPDATE requisitions
SET reasons    = $1,
    status     = $2,
    data       = $3,
    pubkey     = (
        case
            when $5 = '' OR $5 is null then pubkey
            else $5 end
        ),
    updated_at = NOW()
WHERE requisition = $4
`

type UpdateRequisitionParams struct {
	Reasons     string            `json:"reasons"`
	Status      RequisitionStatus `json:"status"`
	Data        []byte            `json:"data"`
	Requisition string            `json:"requisition"`
	Pubkey      interface{}       `json:"pubkey"`
}

func (q *Queries) UpdateRequisition(ctx context.Context, arg UpdateRequisitionParams) error {
	_, err := q.db.ExecContext(ctx, updateRequisition,
		arg.Reasons,
		arg.Status,
		arg.Data,
		arg.Requisition,
		arg.Pubkey,
	)
	return err
}
