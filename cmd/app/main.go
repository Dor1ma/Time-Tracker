package main

import (
	"github.com/Dor1ma/Time-Tracker/internal/app"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	app.Start()
}
