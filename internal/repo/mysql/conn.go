package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mayankpahwa/aspire-loan-app/app/config"
)

var conn *sql.DB

func InitDatabase(config config.Config) error {
	if conn != nil {
		return nil
	}
	db, err := sql.Open("mysql", config.DB.DSN)
	if err != nil {
		return err
	}
	conn = db
	return nil
}

func GetConnection() *sql.DB {
	return conn
}
