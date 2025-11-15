
package handlers

import (
	"anistream/internal/db"
	anistream_minio "anistream/internal/minio"
	"anistream/internal/models"
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

const posterBucketName = "posters"

// AdminCreateAnime adds a new anime to the database.
func AdminCreateAnime(c *gin.Context) {
	var anime models.Anime
	if err := c.ShouldBindJSON(&anime); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Create(&anime).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, anime)
}

// AdminUpdateAnime updates an existing anime in the database.
func AdminUpdateAnime(c *gin.Context) {
	id := c.Param("id")
	var anime models.Anime
	if err := db.DB.First(&anime, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Anime not found"})
		return
	}

	if err := c.ShouldBindJSON(&anime); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Save(&anime).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, anime)
}

// AdminDeleteAnime deletes an anime from the database.
func AdminDeleteAnime(c *gin.Context) {
	id := c.Param("id")
	var anime models.Anime
	if err := db.DB.First(&anime, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Anime not found"})
		return
	}

	if err := db.DB.Delete(&anime).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Anime deleted successfully"})
}

// UploadPoster handles the upload of a poster image to MinIO.
func UploadPoster(c *gin.Context) {
	animeIDStr := c.Param("id")
	animeID, err := strconv.ParseUint(animeIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid anime ID"})
		return
	}

	// 0. Check if anime exists
	var anime models.Anime
	if err := db.DB.First(&anime, animeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Anime not found"})
		return
	}

	// 1. Get the uploaded file
	file, err := c.FormFile("poster")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Poster file not provided"})
		return
	}

	// 2. Create the posters bucket if it doesn't exist
	ctx := context.Background()
	exists, err := anistream_minio.Client.BucketExists(ctx, posterBucketName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check for bucket"})
		return
	}
	if !exists {
		err = anistream_minio.Client.MakeBucket(ctx, posterBucketName, minio.MakeBucketOptions{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create bucket"})
			return
		}
	}

	// 3. Upload the file to MinIO
	fileName := filepath.Base(file.Filename)
	objectName := fmt.Sprintf("%d/%s", animeID, fileName)
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open uploaded file"})
		return
	}
	defer src.Close()

	_, err = anistream_minio.Client.PutObject(ctx, posterBucketName, objectName, src, file.Size, minio.PutObjectOptions{
		ContentType: file.Header.Get("Content-Type"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file to MinIO"})
		return
	}

	// 4. Update the anime record with the new poster URL
	posterURL := fmt.Sprintf("/posters/%s", objectName)
	if err := db.DB.Model(&anime).Update("poster_url", posterURL).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update anime poster URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Poster uploaded successfully", "poster_url": posterURL})
}

