package database

import (
	"fmt"
	"log"
	"sync"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
	_ "github.com/lib/pq"
)

var once sync.Once

// InitializeDB ensures that the database is set up only once
func InitializeDB() {
	once.Do(func() {
		// Read database config from app.conf
		dbDriver, _ := web.AppConfig.String("db_driver")
		dbUser, _ := web.AppConfig.String("db_user")
		dbPassword, _ := web.AppConfig.String("db_password")
		dbName, _ := web.AppConfig.String("db_name")
		dbHost, _ := web.AppConfig.String("db_host")
		dbPort, _ := web.AppConfig.String("db_port")
		dbSSLMode, _ := web.AppConfig.String("db_sslmode")

		// Construct a correct DSN with the created database
		dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
			dbUser, dbPassword, dbName, dbHost, dbPort, dbSSLMode)

		// Register database
		orm.RegisterDriver(dbDriver, orm.DRPostgres)
		orm.RegisterDataBase("default", dbDriver, dsn)

		// Run migrations
		if err := orm.RunSyncdb("default", false, true); err != nil {
			log.Fatalf("Failed to sync database: %v", err)
		}

		log.Println("Database connected and migrated successfully!")
	})
}
