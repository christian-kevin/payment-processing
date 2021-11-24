package mysql

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"spenmo/payment-processing/payment-processing/internal/pkg/component"
	dto "spenmo/payment-processing/payment-processing/internal/pkg/store"
	"spenmo/payment-processing/payment-processing/internal/pkg/timeutil"
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
			created_at
		) values (
			:wallet_id,
			:amount,
			:created_at
		)
	`
)

func (w *walletBalanceLogStore) CreateLog(ctx context.Context, execer Execer, log *dto.WalletBalanceLog) error {
	log.CreatedAt = timeutil.NowMillis()
	_, err := execer.NamedExecContext(ctx, insertWalletBalanceLog, log)
	if err != nil {
		return fmt.Errorf("failed to insert wallet balance log %v: %w", log, err)
	}

	return nil
}

