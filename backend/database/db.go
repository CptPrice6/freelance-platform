package database

import (
	"fmt"
	"log"
	"sync"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

var once sync.Once

func InitializeDB() {
	once.Do(func() {
		dbDriver, _ := web.AppConfig.String("db_driver")
		dbUser, _ := web.AppConfig.String("db_user")
		dbPassword, _ := web.AppConfig.String("db_password")
		dbName, _ := web.AppConfig.String("db_name")
		dbHost, _ := web.AppConfig.String("db_host")
		dbPort, _ := web.AppConfig.String("db_port")
		dbSSLMode, _ := web.AppConfig.String("db_sslmode")

		dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
			dbUser, dbPassword, dbName, dbHost, dbPort, dbSSLMode)

		orm.RegisterDriver(dbDriver, orm.DRPostgres)
		orm.RegisterDataBase("default", dbDriver, dsn)

		if err := orm.RunSyncdb("default", false, true); err != nil {
			log.Fatalf("Failed to sync database: %v", err)
		}

		db, err := goose.OpenDBWithDriver("postgres", dsn)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		if err := goose.Up(db, "migrations"); err != nil {
			log.Fatal(err)
		}

		log.Println("Database connected and migrated successfully!")
	})
}
