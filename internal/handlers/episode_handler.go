
package handlers

import (
	"anistream/internal/db"
	anistream_minio "anistream/internal/minio"
	"anistream/internal/models"
	"anistream/internal/rabbitmq"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

// CreateEpisode adds a new episode to the database
func CreateEpisode(c *gin.Context) {
	var episode models.Episode
	if err := c.ShouldBindJSON(&episode); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Create(&episode).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, episode)
}

// GetEpisode retrieves a single episode from the database
func GetEpisode(c *gin.Context) {
	var episode models.Episode
	if err := db.DB.First(&episode, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Episode not found"})
		return
	}
	c.JSON(http.StatusOK, episode)
}

const rawVideoBucketName = "raw-videos"

// UploadRawVideo handles the upload of a raw video file to MinIO,
// creates a new VideoVersion record, and sends a transcoding job to RabbitMQ.
func UploadRawVideo(c *gin.Context) {
	episodeIDStr := c.Param("id")
	episodeID, err := strconv.ParseUint(episodeIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid episode ID"})
		return
	}

	// 0. Check if episode exists
	var episode models.Episode
	if err := db.DB.First(&episode, episodeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Episode not found"})
		return
	}

	// 1. Get the uploaded file
	file, err := c.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Video file not provided"})
		return
	}

	// 2. Create the raw-videos bucket if it doesn't exist
	ctx := context.Background()
	exists, err := anistream_minio.Client.BucketExists(ctx, rawVideoBucketName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check for bucket"})
		return
	}
	if !exists {
		err = anistream_minio.Client.MakeBucket(ctx, rawVideoBucketName, minio.MakeBucketOptions{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create bucket"})
			return
		}
	}

	// 3. Upload the file to MinIO
	objectName := fmt.Sprintf("%d-%s", episodeID, filepath.Base(file.Filename))
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open uploaded file"})
		return
	}
	defer src.Close()

	_, err = anistream_minio.Client.PutObject(ctx, rawVideoBucketName, objectName, src, file.Size, minio.PutObjectOptions{
		ContentType: file.Header.Get("Content-Type"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file to MinIO"})
		return
	}

	// 4. Create a new VideoVersion record
	videoVersion := models.VideoVersion{
		EpisodeID:  uint(episodeID),
		Quality:    "raw",
		Format:     "mp4", // Assuming mp4 for now
		ObjectPath: objectName,
		Status:     "pending",
	}

	if err := db.DB.Create(&videoVersion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create video version record"})
		return
	}

	// 5. Send a message to RabbitMQ
	job := map[string]uint{"video_version_id": videoVersion.ID}
	body, err := json.Marshal(job)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create JSON for RabbitMQ"})
		return
	}

	if err := rabbitmq.Publish("video_transcode_queue", body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish message to RabbitMQ"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded and transcoding job started successfully"})
}
