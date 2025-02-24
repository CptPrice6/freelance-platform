package main

import (
	"backend/database"
	_ "backend/routers"
	"backend/seeder"

	"github.com/beego/beego/v2/server/web"
)

func main() {
	// Initialize database once
	database.InitializeDB()

	// Seed database with initial data
	seeder.SeedDatabase()

	// Start the Beego application
	web.Run()
}
