package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

const (
	HOST     = "localhost"
	PORT     = 5432
	USER     = "misha"
	PASSWORD = ""
	DBNAME   = "testdb"
)

// Opens a db and returns one and error if not nil
func OpenDb() (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", USER, PASSWORD, HOST, PORT, DBNAME))
	if err != nil {
		return nil, err
	}
	return db, nil
}
