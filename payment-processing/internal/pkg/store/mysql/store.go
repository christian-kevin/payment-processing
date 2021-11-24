package mysql

import (
	"github.com/jmoiron/sqlx"
	"spenmo/payment-processing/payment-processing/internal/pkg/component"
)

type store struct {
	dbx *sqlx.DB
}

func NewStore(db *component.Database) Store {
	return &store{
		dbx: sqlx.NewDb(db.Master, "mysql"),
	}
}
func (s *store) DBX() *sqlx.DB {
	return s.dbx
}

func (s *store) BeginX() (tx *sqlx.Tx, err error) {
	return s.dbx.Beginx()
}

func (s *store) CommitX(tx *sqlx.Tx) error {
	return tx.Commit()
}

func (s *store) RollbackX(tx *sqlx.Tx) error {
	return tx.Rollback()
}
