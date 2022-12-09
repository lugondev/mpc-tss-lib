// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0

package db_gateway

import (
	"context"
)

type Querier interface {
	FailRequisition(ctx context.Context, arg FailRequisitionParams) error
	GetRequisition(ctx context.Context, requisition string) (Requisition, error)
	GetRequisitionById(ctx context.Context, id int64) (Requisition, error)
	GetRetryTimes(ctx context.Context, requisition string) (int32, error)
	InsertRequisition(ctx context.Context, arg InsertRequisitionParams) (Requisition, error)
	ListRequisitions(ctx context.Context, arg ListRequisitionsParams) ([]Requisition, error)
	RetryRequisition(ctx context.Context, requisition string) error
	UpdateRequisition(ctx context.Context, arg UpdateRequisitionParams) error
}

var _ Querier = (*Queries)(nil)