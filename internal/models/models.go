
package models

import "time"

// Anime represents the anime table in the database
type Anime struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	Title            string    `gorm:"not null" json:"title"`
	Description      string    `json:"description"`
	ReleaseYear      int       `json:"release_year"`
	PosterURL        string    `json:"poster_url"`
	Genres           string    `json:"genres"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	Episodes         []Episode `gorm:"foreignKey:AnimeID" json:"episodes"`
}

// Episode represents the episodes table in the database
type Episode struct {
	ID                  uint           `gorm:"primaryKey" json:"id"`
	AnimeID             uint           `gorm:"not null" json:"anime_id"`
	EpisodeNumber       int            `gorm:"not null" json:"episode_number"`
	Title               string         `json:"title"`
	DurationSeconds     int            `json:"duration_seconds"`
	ThumbnailObjectPath string         `json:"thumbnail_object_path"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	VideoVersions       []VideoVersion `gorm:"foreignKey:EpisodeID" json:"video_versions"`
}

// VideoVersion represents the video_versions table in the database
type VideoVersion struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	EpisodeID  uint      `gorm:"not null" json:"episode_id"`
	Quality    string    `gorm:"not null" json:"quality"`
	Format     string    `gorm:"not null" json:"format"`
	ObjectPath string    `gorm:"not null" json:"object_path"`
	Status     string    `gorm:"not null" json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
