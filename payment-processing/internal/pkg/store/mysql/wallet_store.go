package mysql

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"spenmo/payment-processing/payment-processing/internal/pkg/component"
	dto "spenmo/payment-processing/payment-processing/internal/pkg/store"
)

type walletStore struct {
	Store
	dbx *sqlx.DB
}

func NewWalletStore(db *component.Database) WalletStore {
	return &walletStore{
		Store: NewStore(db),
		dbx:   sqlx.NewDb(db.Master, "mysql"),
	}
}

const (
	insertWalletQuery = `
		insert into wallet(
			balance,
			parent_key,
			parent_type,
			country
		) values (
			:balance,
			:parent_key,
			:parent_type,
			:country
		)
	`
)

func (w *walletStore) CreateWallet(ctx context.Context, execer Execer, wallet *dto.Wallet) (walletID int64, err error) {
	res, err := execer.NamedExecContext(ctx, insertWalletQuery, wallet)
	if err != nil {
		return -1, fmt.Errorf("failed to insert wallet %v: %w", wallet, err)
	}

	id, _ := res.LastInsertId()

	return id, nil
}

func (w *walletStore) ModifyBalance(ctx context.Context, execer Execer, inflatedAmount int64, walletID int64) error {
	return nil
}
