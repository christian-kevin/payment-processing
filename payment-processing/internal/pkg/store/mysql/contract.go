package mysql

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
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
	DBX(isMaster bool) *sqlx.DB
	BeginX() (tx *sqlx.Tx, err error)
	CommitX(tx *sqlx.Tx) error
	RollbackX(tx *sqlx.Tx) error
}
