// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0

package db_client

import (
	"context"
)

type Querier interface {
	AddChain(ctx context.Context, arg AddChainParams) (Chain, error)
	AddEmail(ctx context.Context, arg AddEmailParams) (Email, error)
	CreateShare(ctx context.Context, arg CreateShareParams) (CreateShareRow, error)
	DeleteEmail(ctx context.Context, id int64) error
	GetChain(ctx context.Context, chainID int64) (Chain, error)
	GetEmail(ctx context.Context, id int64) (Email, error)
	GetPartyIdByPubkey(ctx context.Context, pubkey string) (GetPartyIdByPubkeyRow, error)
	GetShare(ctx context.Context, id int64) (GetShareRow, error)
	GetShareByAddress(ctx context.Context, address string) (GetShareByAddressRow, error)
	GetShareByID(ctx context.Context, partyID string) (GetShareByIDRow, error)
	ListChains(ctx context.Context) ([]Chain, error)
	ListEmails(ctx context.Context, arg ListEmailsParams) ([]Email, error)
	ListShare(ctx context.Context) ([]ListShareRow, error)
	UpdateEmail(ctx context.Context, arg UpdateEmailParams) (Email, error)
}

var _ Querier = (*Queries)(nil)