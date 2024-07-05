package app

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	migratePg "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/sirupsen/logrus"
)

func RunMigration(dsn string, log *logrus.Logger) {
	dbMigration, err := sql.Open("postgres", dsn)
	defer dbMigration.Close()

	if err != nil {
		log.Debugf("Failed to connect to database: %v", err)
		return
	}

	driver, err := migratePg.WithInstance(dbMigration, &migratePg.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return
	}

	err = m.Up()
	if err != nil {
		log.Debugf("Failed to run migration: %v", err)
	}

	log.Info("Migrations ran successfully")
}
