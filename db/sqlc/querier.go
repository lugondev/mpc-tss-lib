// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0

package db

import (
	"context"
)

type Querier interface {
	AddChain(ctx context.Context, arg AddChainParams) (Chain, error)
	AddEmail(ctx context.Context, arg AddEmailParams) (Email, error)
	CreateShare(ctx context.Context, arg CreateShareParams) (Share, error)
	DeleteEmail(ctx context.Context, id int64) error
	GetChain(ctx context.Context, chainID int64) (Chain, error)
	GetEmail(ctx context.Context, id int64) (Email, error)
	GetShare(ctx context.Context, id int64) (Share, error)
	GetShareByAddress(ctx context.Context, address string) (Share, error)
	GetShareByPubkey(ctx context.Context, pubkey string) (Share, error)
	ListChains(ctx context.Context) ([]Chain, error)
	ListEmails(ctx context.Context, arg ListEmailsParams) ([]Email, error)
	ListShare(ctx context.Context, arg ListShareParams) ([]Share, error)
	UpdateEmail(ctx context.Context, arg UpdateEmailParams) (Email, error)
}

var _ Querier = (*Queries)(nil)
