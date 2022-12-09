package db_client

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lugondev/mpc-tss-lib/db/client/sqlc"
	"github.com/rs/zerolog"
)

// Store defines all functions to execute db queries and transactions
type Store interface {
	db_client.Querier
}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	db *sql.DB
	*db_client.Queries

	logger zerolog.Logger
}

// NewStore creates a new store
func NewStore(sqlDb *sql.DB, logger *zerolog.Logger) *SQLStore {
	connLogger := logger.With().Str("clientID", "db-gateway").Logger()
	return &SQLStore{
		db:      sqlDb,
		Queries: db_client.New(sqlDb),
		logger:  connLogger,
	}
}

// ExecTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*db_client.Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := db_client.New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
