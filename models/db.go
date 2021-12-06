package models

import (
	"example/hello/config"

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

	db, err := gorm.Open(postgres.Open(config.DB.URI), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: config.DB.Prefix,
		},
	})

	if err != nil {
		panic("failed to connect database")
	}

	return db
}
