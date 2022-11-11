package main

import (
	"os"
	"log"
	"fmt"
	"testing"
	"database/sql"
	_ "github.com/lib/pq"
)


var (
	dbManager *DBManager
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "1234"
	dbname   = "shop"

)

func NewDBmanager(db *sql.DB) *DBManager {
	return &DBManager{db: db}
}

func TestMain(m *testing.M) {
	condtb := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", condtb)
	if err != nil {
		log.Fatalf("Error while opening database: %v", err)
	}

	dbManager = NewDBmanager(db)
	os.Exit(m.Run())
}