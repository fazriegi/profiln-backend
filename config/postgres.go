package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func NewDatabase() *sql.DB {
	username := os.Getenv("DBUSER")
	password := os.Getenv("DBPASS")
	host := os.Getenv("DBHOST")
	port := os.Getenv("DBPORT")
	name := os.Getenv("DBNAME")

	dbDriver := "postgres"
	dbSource := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
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

	// !! uncomment this if you want to automate migration without installing migrate cli
	// driver, err := postgres.WithInstance(conn, &postgres.Config{})
	// if err != nil {
	// 	log.Fatal("postgres.WithInstance", err)
	// }

	// m, err := migrate.NewWithDatabaseInstance(
	// 	"file://db/migrations",
	// 	"postgres", driver)
	// if err != nil {
	// 	log.Fatal("migrate.NewWithDatabaseInstance", err)
	// }

	// m.Up()

	return conn
}
