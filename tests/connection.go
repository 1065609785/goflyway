package tests

import (
	_ "com.goldstar/go-sql-driver/mysql"
	_ "com.goldstar/mattn/go-sqlite3"
	"database/sql"
)

func ConnSqlite() (*sql.DB, error) {
	return sql.Open("sqlite3", "./flyway_test.db")
}

func ConnectMysql() (*sql.DB, error) {
	return sql.Open("mysql", "root:goflyway@tcp(localhost:9910)/goflyway?charset=utf8")
}
