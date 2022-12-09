package db_gateway

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/lugondev/mpc-tss-lib/db/gateway/sqlc"
	"github.com/rs/zerolog"
)

// Store defines all functions to execute db queries and transactions
type Store interface {
	db_gateway.Querier
}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	db *sql.DB
	*db_gateway.Queries

	logger zerolog.Logger
}

// NewStore creates a new store
func NewStore(sqlDb *sql.DB, logger *zerolog.Logger) *SQLStore {
	connLogger := logger.With().Str("clientID", "db-gateway").Logger()
	return &SQLStore{
		db:      sqlDb,
		Queries: db_gateway.New(sqlDb),
		logger:  connLogger,
	}
}

type RequisitionType string

const (
	RequisitionSign    RequisitionType = "sign"
	RequisitionKeygen  RequisitionType = "keygen"
	RequisitionReshare RequisitionType = "reshare"
)

func (store *SQLStore) CreateRequisition(ctx context.Context, params db_gateway.InsertRequisitionParams, requisitionType RequisitionType) (db_gateway.Requisition, error) {
	store.logger.Info().Msgf("Creating requisition for %s ", requisitionType)
	store.logger.Info().Msgf("pubkey %s ", params.Pubkey)

	params.Status = "pending"
	params.Requisition = fmt.Sprintf("%s", uuid.New())
	if requisitionType != RequisitionKeygen && params.Pubkey == "" {
		return db_gateway.Requisition{}, fmt.Errorf("pubkey is required")
	}

	switch requisitionType {

	case RequisitionSign:
		params.Type = "sign"

	case RequisitionKeygen:
		params.Type = "keygen"

	case RequisitionReshare:
		params.Type = "reshare"

	}

	return store.InsertRequisition(ctx, params)
}
