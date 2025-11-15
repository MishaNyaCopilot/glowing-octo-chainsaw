
package db

import (
	"anistream/internal/models"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Init initializes the database connection and performs auto-migration
func Init(dsn string) {
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate the schema
	DB.AutoMigrate(&models.Anime{}, &models.Episode{}, &models.VideoVersion{})
	fmt.Println("Database migration completed!")
}
