
package main

import (
	"anistream/internal/config"
	"anistream/internal/db"
	"anistream/internal/minio"
	"anistream/internal/rabbitmq"
	"anistream/internal/router"
	"fmt"
	"log"
)

func main() {
	cfg := config.New()

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db.Init(dsn)

	rabbitmq.Init(cfg.RabbitMQUri)
	defer rabbitmq.Close()

	minio.NewClient(cfg.MinioHost, cfg.MinioRootUser, cfg.MinioRootPass, false)

	r := router.SetupRouter()

	log.Println("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
