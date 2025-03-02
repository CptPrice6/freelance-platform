package main

import (
	"backend/database"
	_ "backend/routers"
	"backend/seeder"

	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/filter/cors"
)

func main() {
	// Initialize database once
	database.InitializeDB()

	// Seed database with initial data
	seeder.SeedDatabase()

	// Initialize CORS
	web.InsertFilter("*", web.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Start the Beego application
	web.Run()
}
