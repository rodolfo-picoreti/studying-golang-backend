package models

import (
	"github.com/rodolfo-picoreti/studying-golang-backend/config"
	"github.com/rodolfo-picoreti/studying-golang-backend/telemetry"

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
	db.AutoMigrate(&Attribute{})
	db.AutoMigrate(&AttributeDefinition{})
	db.AutoMigrate(&ProductType{})
	db.AutoMigrate(&Product{})
	db.AutoMigrate(&Inventory{})
	db.AutoMigrate(&Price{})
}
