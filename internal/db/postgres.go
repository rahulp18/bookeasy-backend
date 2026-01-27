package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	connStr := "postgres://postgres:root@localhost:5432/bookeasy?sslmode=disable"
	var err error
	DB, err = sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal("failed to open db", err)
	}
	err = DB.Ping()
	if err != nil {
		log.Fatal("failed to ping db:", err)
	}
	fmt.Println("PostgresSQL connected successfully")
}
