package app

import (
	"fmt"
	"github.com/Dor1ma/Time-Tracker/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var db *gorm.DB

func Start() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DbHost,
		cfg.DbUser,
		cfg.DbPass,
		cfg.DbName,
		cfg.DbPort)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	RunMigration(dsn)
}
