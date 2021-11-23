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
	Store
	CreateWallet(ctx context.Context, execer Execer, wallet *dto.Wallet) (walletID int64, err error)
	ModifyBalance(ctx context.Context, execer Execer, inflatedNewBalance int64, walletID int64) error
	GetWalletByUserID(ctx context.Context, querier Querier, userID int64) (*dto.Wallet, error)
}

type CardStore interface {
	Store
	CreateCard(ctx context.Context, execer Execer, card *dto.Card) (cardID int64, err error)
	GetCardByID(ctx context.Context, querier Querier, cardID int64) (*dto.Card, error)
	GetCards(ctx context.Context, querier Querier, walletID int64) ([]*dto.Card, error)
	GetCardByNumberAndExpiryDate(ctx context.Context, querier Querier, cardNumber string,
		expiryDate string) (*dto.Card, error)
	DeleteCard(ctx context.Context, execer Execer, cardID int64) error
}

type CardTransactionLogStore interface {
	Store
	CreateLog(ctx context.Context, execer Execer, log *dto.CardTransactionLog) error
}

type WalletBalanceLogStore interface {
	Store
	CreateLog(ctx context.Context, execer Execer, log *dto.WalletBalanceLog) error
}

type LimitStore interface {
	Store
	CreateLimit(ctx context.Context, execer Execer, limit *dto.Limit) (limitID int64, err error)
	GetLimits(ctx context.Context, querier Querier, parentType int, parentID int64) ([]*dto.Limit, error)
}
