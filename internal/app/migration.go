package app

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	migratePg "github.com/golang-migrate/migrate/v4/database/postgres"
	"log"
)

func RunMigration(dsn string) {
	dbMigration, err := sql.Open("postgres", dsn)
	defer dbMigration.Close()

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	driver, err := migratePg.WithInstance(dbMigration, &migratePg.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	m.Up()

	log.Println("Migrations ran successfully")
}
