
package router

import (
	"anistream/internal/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures the Gin router with all the application routes
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Add CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"}
	r.Use(cors.New(config))

	api := r.Group("/api")
	{
		// Public routes
		anime := api.Group("/anime")
		{
			anime.GET("/", handlers.GetAnimes)
			anime.GET("/:id", handlers.GetAnime)
		}

		episode := api.Group("/episode")
		{
			episode.GET("/:id", handlers.GetEpisode)
		}

		// Admin routes
		admin := api.Group("/admin")
		{
			// Anime routes
			admin.POST("/anime", handlers.AdminCreateAnime)
			admin.PUT("/anime/:id", handlers.AdminUpdateAnime)
			admin.DELETE("/anime/:id", handlers.AdminDeleteAnime)
			admin.POST("/anime/:id/upload-poster", handlers.UploadPoster)

			// Episode routes
			admin.POST("/episode", handlers.AdminCreateEpisode)
			admin.PUT("/episode/:id", handlers.AdminUpdateEpisode)
			admin.DELETE("/episode/:id", handlers.AdminDeleteEpisode)
			admin.POST("/episode/:id/upload-raw-video", handlers.UploadRawVideo)
		}
	}

	hls := r.Group("/hls")
	{
		hls.GET("/:episodeID/playlist.m3u8", handlers.GetHLSPlaylist)
		hls.GET("/segments/:episodeID/:segmentFile", handlers.GetHLSSegment)
	}

	posters := r.Group("/posters")
	{
		posters.GET("/:animeID/:posterFile", handlers.GetPoster)
	}

	return r
}
