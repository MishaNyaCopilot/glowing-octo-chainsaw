
package handlers

import (
	"anistream/internal/db"
	"anistream/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminCreateEpisode adds a new episode to the database.
func AdminCreateEpisode(c *gin.Context) {
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

// AdminUpdateEpisode updates an existing episode in the database.
func AdminUpdateEpisode(c *gin.Context) {
	id := c.Param("id")
	var episode models.Episode
	if err := db.DB.First(&episode, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Episode not found"})
		return
	}

	if err := c.ShouldBindJSON(&episode); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Save(&episode).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, episode)
}

// AdminDeleteEpisode deletes an episode from the database.
func AdminDeleteEpisode(c *gin.Context) {
	id := c.Param("id")
	var episode models.Episode
	if err := db.DB.First(&episode, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Episode not found"})
		return
	}

	if err := db.DB.Delete(&episode).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Episode deleted successfully"})
}
