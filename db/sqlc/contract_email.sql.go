// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: contract_email.sql

package db

import (
	"context"
	"time"
)

const listEmailsSubscription = `-- name: ListEmailsSubscription :many
SELECT emails.id, emails.name, emails.email, emails.created_at
FROM emails_contract
         JOIN emails ON emails_contract.email_id = emails.id
WHERE contract_id = $1
`

func (q *Queries) ListEmailsSubscription(ctx context.Context, contractid int64) ([]Email, error) {
	rows, err := q.db.QueryContext(ctx, listEmailsSubscription, contractid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Email{}
	for rows.Next() {
		var i Email
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Email,
			&i.CreatedAt,
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

const listEmailsSubscriptionByAddress = `-- name: ListEmailsSubscriptionByAddress :many
SELECT emails.id, emails.name, emails.email, emails.created_at
FROM contracts
         JOIN emails_contract ON emails_contract.contract_id = contracts.id
         JOIN emails ON emails_contract.email_id = emails.id
WHERE contracts.address = $1
`

func (q *Queries) ListEmailsSubscriptionByAddress(ctx context.Context, contractaddress string) ([]Email, error) {
	rows, err := q.db.QueryContext(ctx, listEmailsSubscriptionByAddress, contractaddress)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Email{}
	for rows.Next() {
		var i Email
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Email,
			&i.CreatedAt,
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

const mapEmailContract = `-- name: MapEmailContract :one
WITH inserted_data AS (
    INSERT
        INTO emails_contract
            (email_id, contract_id)
            VALUES ($1, $2) RETURNING id, email_id, contract_id, created_at)
SELECT inserted_data.id, email_id, contract_id, inserted_data.created_at, emails.id, emails.name, email, emails.created_at, contracts.id, contracts.name, is_contract, chain_id, notification, address, network, contracts.created_at
FROM inserted_data
         INNER JOIN emails ON inserted_data.email_id = emails.id
         INNER JOIN contracts ON inserted_data.contract_id = contracts.id
WHERE contracts.id = inserted_data.contract_id
  AND emails.id = inserted_data.email_id
`

type MapEmailContractParams struct {
	EmailID    int64 `json:"email_id"`
	ContractID int64 `json:"contract_id"`
}

type MapEmailContractRow struct {
	ID           int64              `json:"id"`
	EmailID      int64              `json:"email_id"`
	ContractID   int64              `json:"contract_id"`
	CreatedAt    time.Time          `json:"created_at"`
	ID_2         int64              `json:"id_2"`
	Name         string             `json:"name"`
	Email        string             `json:"email"`
	CreatedAt_2  time.Time          `json:"created_at_2"`
	ID_3         int64              `json:"id_3"`
	Name_2       string             `json:"name_2"`
	IsContract   bool               `json:"is_contract"`
	ChainID      string             `json:"chain_id"`
	Notification NotificationStatus `json:"notification"`
	Address      string             `json:"address"`
	Network      string             `json:"network"`
	CreatedAt_3  time.Time          `json:"created_at_3"`
}

func (q *Queries) MapEmailContract(ctx context.Context, arg MapEmailContractParams) (MapEmailContractRow, error) {
	row := q.db.QueryRowContext(ctx, mapEmailContract, arg.EmailID, arg.ContractID)
	var i MapEmailContractRow
	err := row.Scan(
		&i.ID,
		&i.EmailID,
		&i.ContractID,
		&i.CreatedAt,
		&i.ID_2,
		&i.Name,
		&i.Email,
		&i.CreatedAt_2,
		&i.ID_3,
		&i.Name_2,
		&i.IsContract,
		&i.ChainID,
		&i.Notification,
		&i.Address,
		&i.Network,
		&i.CreatedAt_3,
	)
	return i, err
}