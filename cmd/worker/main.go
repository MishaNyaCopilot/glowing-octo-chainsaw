
package main

import (
	"anistream/internal/config"
	"anistream/internal/db"
	"anistream/internal/minio"
	"anistream/internal/rabbitmq"
	"anistream/internal/worker"
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

// TranscodeJob represents the message structure for a transcoding job
type TranscodeJob struct {
	VideoVersionID uint `json:"video_version_id"`
}

func main() {
	cfg := config.New()

	// Initialize DB
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	db.Init(dsn)

	// Initialize MinIO
	minio.NewClient(cfg.MinioHost, cfg.MinioRootUser, cfg.MinioRootPass, false)

	// Initialize RabbitMQ
	rabbitmq.Init(cfg.RabbitMQUri)
	defer rabbitmq.Close()

	msgs, err := rabbitmq.Consume("video_transcode_queue")
	if err != nil {
		log.Fatalf("Failed to register a consumer: %s", err)
	}

	// Number of concurrent workers
	numWorkers := 5
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			log.Printf("Worker %d started", workerID)
			for d := range msgs {
				log.Printf("Worker %d received a message: %s", workerID, d.Body)
				var job TranscodeJob
				if err := json.Unmarshal(d.Body, &job); err != nil {
					log.Printf("Worker %d: Error decoding JSON: %s", workerID, err)
					continue
				}

				if err := worker.TranscodeVideo(job.VideoVersionID); err != nil {
					log.Printf("Worker %d: Error transcoding video: %s", workerID, err)
					// Here you might want to update the video version status to 'failed'
				}
			}
			log.Printf("Worker %d finished", workerID)
		}(i + 1)
	}

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	wg.Wait()
}
