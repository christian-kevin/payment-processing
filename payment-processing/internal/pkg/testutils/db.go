package testutils

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

// PrepareMySQLMock prepares sqlx.DB and Sqlmock with MySQL driver
func PrepareMySQLMock(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	return PrepareDBMock(t, "mysql")
}

// PrepareDBMock prepares sqlx.DB and Sqlmock given database driver
func PrepareDBMock(t *testing.T, driver string) (*sqlx.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	return sqlx.NewDb(db, driver), mock
}
