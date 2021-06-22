package mysql

import (
	"github.com/sendinblue/dpe-insights/core/config"
	_ "github.com/go-sql-driver/mysql" // required for mysql driver
	"github.com/jmoiron/sqlx"
)

// NewDB new mysql databases.
func NewDB(appConfig *config.Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("mysql", appConfig.MysqlDSN)
	if err != nil {
		return nil, err
	}

	return db, err
}
