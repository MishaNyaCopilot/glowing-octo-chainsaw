package worker

import (
	anistream_db "anistream/internal/db"
	anistream_minio "anistream/internal/minio"
	"anistream/internal/models"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

const rawVideoBucket = "raw-videos"
const hlsBucket = "anistream"

// isNvidiaEncoderAvailable checks if the h264_nvenc encoder is available in ffmpeg.
func isNvidiaEncoderAvailable() bool {
	cmd := exec.Command("ffmpeg", "-encoders")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), "h264_nvenc")
}

// getEncoder returns the best available encoder.
func getEncoder() string {
	if isNvidiaEncoderAvailable() {
		log.Println("NVIDIA encoder found, using h264_nvenc.")
		return "h264_nvenc"
	}
	log.Println("NVIDIA encoder not found, falling back to libx264.")
	return "libx264"
}

// TranscodeVideo downloads a raw video, transcodes it to HLS, and uploads the results.
func TranscodeVideo(videoVersionID uint) error {
	log.Printf("Starting transcoding for video version ID: %d", videoVersionID)

	// 1. Get VideoVersion from DB
	var videoVersion models.VideoVersion
	if err := anistream_db.DB.First(&videoVersion, videoVersionID).Error; err != nil {
		return fmt.Errorf("failed to get video version from db: %w", err)
	}

	// 2. Download Raw Video from MinIO
	rawVideoPath := filepath.Join(os.TempDir(), filepath.Base(videoVersion.ObjectPath))
	log.Printf("Downloading raw video from %s to %s", videoVersion.ObjectPath, rawVideoPath)
	err := anistream_minio.Client.FGetObject(context.Background(), rawVideoBucket, videoVersion.ObjectPath, rawVideoPath, minio.GetObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to download raw video: %w", err)
	}
	defer os.Remove(rawVideoPath) // Clean up the downloaded file

	// 3. Transcode to HLS
	outputDir := filepath.Join(os.TempDir(), fmt.Sprintf("hls_%d", videoVersion.EpisodeID))
	os.MkdirAll(outputDir, os.ModePerm)
	defer os.RemoveAll(outputDir) // Clean up the HLS files

	masterPlaylistPath := filepath.Join(outputDir, "master.m3u8")

	encoder := getEncoder()

	log.Printf("Transcoding video to HLS in %s using %s encoder", outputDir, encoder)
	args := ffmpeg.KwArgs{
		"c:v":           encoder,
		"c:a":           "aac",
		"pix_fmt":       "yuv420p",
		"hls_time":      "10",
		"hls_playlist_type": "vod",
		"hls_base_url": fmt.Sprintf("/hls/segments/%d/", videoVersion.EpisodeID),
		"hls_segment_filename": filepath.Join(outputDir, "segment%03d.ts"),
	}

	if encoder == "h264_nvenc" {
		args["preset"] = "p5"
		args["tune"] = "hq"
	}

	err = ffmpeg.Input(rawVideoPath).
		Output(masterPlaylistPath, args).
		Run()

	if err != nil {
		return fmt.Errorf("failed to transcode video: %w", err)
	}

	// 4. Upload HLS files to MinIO
	log.Printf("Uploading HLS files to MinIO bucket: %s", hlsBucket)
	files, err := os.ReadDir(outputDir)
	if err != nil {
		return fmt.Errorf("failed to read HLS output directory: %w", err)
	}

	for _, file := range files {
		filePath := filepath.Join(outputDir, file.Name())
		objectName := fmt.Sprintf("hls/%d/%s", videoVersion.EpisodeID, file.Name())

		_, err := anistream_minio.Client.FPutObject(context.Background(), hlsBucket, objectName, filePath, minio.PutObjectOptions{})
		if err != nil {
			return fmt.Errorf("failed to upload %s to minio: %w", file.Name(), err)
		}
	}

	// 5. Update VideoVersion in DB
	log.Printf("Updating video version status to 'ready'")
	masterPlaylistObjectName := fmt.Sprintf("hls/%d/master.m3u8", videoVersion.EpisodeID)
	result := anistream_db.DB.Model(&videoVersion).Updates(models.VideoVersion{
		Status:     "ready",
		ObjectPath: masterPlaylistObjectName,
	})
	if result.Error != nil {
		return fmt.Errorf("failed to update video version in db: %w", result.Error)
	}

	log.Printf("Transcoding finished successfully for video version ID: %d", videoVersionID)
	return nil
}