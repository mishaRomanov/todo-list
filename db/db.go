package db

import (
	"database/sql"
	_ "github.com/lib/pq"
)

func OpenDb() *sql.DB {
	connStr := "user=guest password= dbname=ids sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	return db
}
