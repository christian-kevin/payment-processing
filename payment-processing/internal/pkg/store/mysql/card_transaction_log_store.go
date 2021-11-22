package mysql

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"spenmo/payment-processing/payment-processing/internal/pkg/component"
	dto "spenmo/payment-processing/payment-processing/internal/pkg/store"
)

type cardTransactionLogStore struct {
	Store
	dbx *sqlx.DB
}

func NewCardTransactionLogStore(db *component.Database) CardTransactionLogStore {
	return &cardTransactionLogStore{
		Store: NewStore(db),
		dbx:   sqlx.NewDb(db.Master, "mysql"),
	}
}

const (
	insertCardTransactionLog = `
		insert into card_transaction_log (
			card_id,
			amount,
			created_at
		) values (
			:card_id,
			:amount,
			:created_at
		)
	`
)

func (c *cardTransactionLogStore) CreateLog(ctx context.Context, execer Execer, log *dto.CardTransactionLog) error {
	_, err := execer.NamedExecContext(ctx, insertCardTransactionLog, log)
	if err != nil {
		return fmt.Errorf("failed to insert card transaction log %v: %w", log, err)
	}

	return nil
}
