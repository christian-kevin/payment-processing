package mysql

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	dto "spenmo/payment-processing/payment-processing/internal/pkg/store"
)

type Querier interface {
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	sqlx.Ext
}

type Execer interface {
	sqlx.ExecerContext
	sqlx.Ext
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}

type Store interface {
	DBX() *sqlx.DB
	BeginX() (tx *sqlx.Tx, err error)
	CommitX(tx *sqlx.Tx) error
	RollbackX(tx *sqlx.Tx) error
}

type WalletStore interface {
	CreateWallet(ctx context.Context, execer Execer, wallet *dto.Wallet) (walletID int64, err error)
	ModifyBalance(ctx context.Context, execer Execer, inflatedAmount int64, walletID int64) error
}

type CardStore interface {
	CreateCard(ctx context.Context, execer Execer, walletID int64, cardNumber, expiryDate, name string) (cardID int64, err error)
	GetCards(ctx context.Context, querier Querier, walletID int64) ([]*dto.Card, error)
	GetCardDetail(ctx context.Context, querier Querier, cardID int64) (*dto.Card, error)
	DeleteCard(ctx context.Context, execer Execer, cardID int64) error
}

type CardTransactionLog interface {
	CreateLog(ctx context.Context, execer Execer, log *dto.CardTransactionLog) error
}

type WalletBalanceLog interface {
	CreateLog(ctx context.Context, execer Execer, log *dto.WalletBalanceLog) error
}

type Limit interface {
	CreateLimit(ctx context.Context, execer Execer, limit *dto.Limit) (limitID int64, err error)
	GetLimit(ctx context.Context, querier Querier, parentType int, parentID int64) (*dto.Limit, error)
}
