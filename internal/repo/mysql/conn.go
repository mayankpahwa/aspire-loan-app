package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var conn *sql.DB

func InitDatabase() error {
	if conn != nil {
		return nil
	}
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/aspire")
	if err != nil {
		return err
	}
	conn = db
	return nil
}

func GetConnection() *sql.DB {
	if conn != nil {
		InitDatabase()
	}
	return conn
}
