package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/tutor_db?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	sqlDB, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)

	}

	testQueries = New(sqlDB)

	os.Exit(m.Run())
}
