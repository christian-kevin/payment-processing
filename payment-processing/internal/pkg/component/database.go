package component

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"spenmo/payment-processing/payment-processing/config"
)

const (
	DBType             = "mysql"
	DBConnectionString = "%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true"
)

type Database struct {
	Master *sql.DB
	Slave  *sql.DB
}

func InitializeDatabase() (database *Database, err error) {
	database = &Database{}
	// default database
	connectionString := fmt.Sprintf(DBConnectionString,
		config.AppConfig.MySQLUser,
		config.AppConfig.MySQLPass,
		config.AppConfig.MySQLHost,
		config.AppConfig.MySQLPort,
		config.AppConfig.MySQLDB,
	)

	database.Master, err = sql.Open(DBType, connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open default database connection. %+v", err)
	}
	database.Master.SetMaxIdleConns(config.AppConfig.MySQLMaxIdle)
	database.Master.SetMaxOpenConns(config.AppConfig.MySQLMaxOpen)

	err = database.Master.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping default database. %+v", err)
	}

	return
}
