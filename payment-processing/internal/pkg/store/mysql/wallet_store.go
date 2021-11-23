package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"spenmo/payment-processing/payment-processing/internal/pkg/component"
	"spenmo/payment-processing/payment-processing/internal/pkg/constant"
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
			parent_id,
			parent_type,
			country
		) values (
			:balance,
			:parent_id,
			:parent_type,
			:country
		)
	`

	updateWalletBalanceQuery = `
		update wallet set
			balance = ?
		where 
			id = ?
	`

	getWalletByParentIDAndParentKey = `
		select 
			id,
			balance,
			parent_id,
			parent_type,
			country
		from
			wallet
		where
			parent_id = ?  and parent_type = ?
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

func (w *walletStore) ModifyBalance(ctx context.Context, execer Execer, inflatedNewBalance int64, walletID int64) error {
	if _, err := execer.ExecContext(ctx, updateWalletBalanceQuery, inflatedNewBalance, walletID); err != nil {
		return fmt.Errorf("failed to update wallet balance for id %d: %w", walletID, err)
	}

	return nil
}

func (w *walletStore) GetWalletByUserID(ctx context.Context, querier Querier, userID int64) (*dto.Wallet, error) {
	var wallet dto.Wallet
	if err := querier.GetContext(ctx, &wallet, getWalletByParentIDAndParentKey, userID, constant.ParentTypeUser); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to find wallet by user id %d: %w", userID, err)
	}

	return &wallet, nil
}
