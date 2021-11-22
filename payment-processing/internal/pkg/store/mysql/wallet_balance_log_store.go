package mysql

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"spenmo/payment-processing/payment-processing/internal/pkg/component"
	dto "spenmo/payment-processing/payment-processing/internal/pkg/store"
)

type walletBalanceLogStore struct {
	Store
	dbx *sqlx.DB
}

func NewWalletBalanceLogStore(db *component.Database) WalletBalanceLogStore {
	return &walletBalanceLogStore{
		Store: NewStore(db),
		dbx:   sqlx.NewDb(db.Master, "mysql"),
	}
}

const (
	insertWalletBalanceLog = `
		insert into wallet_balance_log (
			wallet_id,
			amount,
			type,
			created_at
		) values (
			:wallet_id,
			:amount,
			:type,
			:created_at
		)
	`
)

func (w *walletBalanceLogStore) CreateLog(ctx context.Context, execer Execer, log *dto.WalletBalanceLog) error {
	_, err := execer.NamedExecContext(ctx, insertWalletBalanceLog, log)
	if err != nil {
		return fmt.Errorf("failed to insert wallet balance log %v: %w", log, err)
	}

	return nil
}

