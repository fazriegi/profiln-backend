package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func NewDatabase() *sql.DB {
	username := os.Getenv("DBUSER")
	password := os.Getenv("DBPASS")
	host := os.Getenv("DBHOST")
	port := os.Getenv("DBPORT")
	name := os.Getenv("DBNAME")

	dbDriver := "postgres"
	dbSource := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		username,
		password,
		host,
		port,
		name,
	)

	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	return conn
}
