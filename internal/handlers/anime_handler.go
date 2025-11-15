
package handlers

import (
	"anistream/internal/db"
	"anistream/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateAnime adds a new anime to the database
func CreateAnime(c *gin.Context) {
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

// GetAnimes retrieves all animes from the database
func GetAnimes(c *gin.Context) {
	var animes []models.Anime
	if err := db.DB.Find(&animes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, animes)
}

// GetAnime retrieves a single anime from the database
func GetAnime(c *gin.Context) {
	var anime models.Anime
	if err := db.DB.Preload("Episodes").First(&anime, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Anime not found"})
		return
	}
	c.JSON(http.StatusOK, anime)
}
