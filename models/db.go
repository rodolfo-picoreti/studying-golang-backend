package models

import (
	"example/hello/config"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var db *gorm.DB

// GetDbConnection returns current database conneciton or creates on the first call
func GetDbConnection() *gorm.DB {
	if db == nil {
		db = createDbConnection()
	}
	return db
}

func createDbConnection() *gorm.DB {
	config := config.ReadConfig()

	log.Println("Connecting to dabatase...")
	db, err := gorm.Open(postgres.Open(config.DB.URI), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: config.DB.Prefix,
		},
	})

	// db.Use(prometheus.New(prometheus.Config{
	// 	DBName:          "db",
	// 	RefreshInterval: 15,
	// 	StartServer:     true,
	// 	HTTPServerPort:  8081,
	// }))

	if err != nil {
		panic("Failed to connect to database")
	}

	return db
}
