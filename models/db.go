package models

import (
	"example/hello/config"
	"example/hello/telemetry"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var db *gorm.DB

// GetDbConnection returns current database conneciton or creates on the first call
func GetDbConnection() *gorm.DB {
	return db
}

func InitDB() {
	config := config.ReadConfig()

	telemetry.GetLogger().Info().Msg("Connecting to dabatase...")
	dbConn, err := gorm.Open(postgres.Open(config.DB.URI), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: config.DB.Prefix,
		},
	})

	if err != nil {
		telemetry.GetLogger().Fatal().Msg("Failed to connect to database")
	}

	db = dbConn

	// db.Use(prometheus.New(prometheus.Config{
	// 	DBName:          "db",
	// 	RefreshInterval: 15,
	// 	StartServer:     true,
	// 	HTTPServerPort:  8081,
	// }))
}

func AutoMigrate() {
	db.AutoMigrate(&Product{})
}
