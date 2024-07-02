package main

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"gorm.io/driver/postgres"
	"log"
	"os"

	migratePg "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // Драйвер PostgreSQL для sql.Open
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"))

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	runMigrations(dsn)
}

func runMigrations(dsn string) {
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
