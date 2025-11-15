package handlers

import (
	"anistream/internal/db"
	anistream_minio "anistream/internal/minio"
	"anistream/internal/models"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

const hlsBucketName = "anistream"

// GetHLSPlaylist serves the HLS master playlist
func GetHLSPlaylist(c *gin.Context) {
	episodeID := c.Param("episodeID")

	var videoVersion models.VideoVersion
	// Find the master playlist for the given episode.
	// We'll assume a 'master' quality exists, or you can pick the highest quality.
	if err := db.DB.Where("episode_id = ? AND status = ?", episodeID, "ready").First(&videoVersion).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "HLS playlist not found for the given episode"})
		return
	}

	// Get the playlist object from MinIO
	object, err := anistream_minio.Client.GetObject(context.Background(), hlsBucketName, videoVersion.ObjectPath, minio.GetObjectOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve playlist from storage"})
		return
	}
	defer object.Close()

	// Check if the object is valid
	if _, err := object.Stat(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve playlist metadata"})
		return
	}

	c.Header("Content-Type", "application/vnd.apple.mpegurl")
	io.Copy(c.Writer, object)
}

// GetHLSSegment serves the HLS video segments
func GetHLSSegment(c *gin.Context) {
	episodeID := c.Param("episodeID")
	segmentFile := c.Param("segmentFile")
	objectPath := fmt.Sprintf("hls/%s/%s", episodeID, segmentFile)

	// Get the segment object from MinIO
	object, err := anistream_minio.Client.GetObject(context.Background(), hlsBucketName, objectPath, minio.GetObjectOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve segment from storage"})
		return
	}
	defer object.Close()

	// Check if the object is valid
	if _, err := object.Stat(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve segment metadata"})
		return
	}

	c.Header("Content-Type", "video/MP2T")
	io.Copy(c.Writer, object)
}

// GetPoster serves the poster image
func GetPoster(c *gin.Context) {
	animeID := c.Param("animeID")
	posterFile := c.Param("posterFile")
	objectPath := fmt.Sprintf("%s/%s", animeID, posterFile)

	// Get the poster object from MinIO
	object, err := anistream_minio.Client.GetObject(context.Background(), posterBucketName, objectPath, minio.GetObjectOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve poster from storage"})
		return
	}
	defer object.Close()

	// Check if the object is valid
	if _, err := object.Stat(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve poster metadata"})
		return
	}

	c.Header("Content-Type", "image/jpeg") // Assuming jpeg, you might want to make this dynamic
	io.Copy(c.Writer, object)
}
