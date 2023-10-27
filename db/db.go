package db

import (
	"database/sql"
	_ "github.com/lib/pq"
)

// Opens a db
func OpenDb() (*sql.DB, error) {
	connStr := "user=guest password= dbname=ids sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	return db, nil
}
