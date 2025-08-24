package ffmpeg

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestFFmpegProcessor(t *testing.T) {
	// 配置媒体播放器
	cfg := Config{
		MinioEndpoint:  "localhost:9000",
		MinioAccessKey: "minioadmin",
		MinioSecretKey: "minioadmin",
		MinioBucket:    "media-bucket",
		MinioUseSSL:    false,
		ServerPort:     "8080",
		FFmpegPath:     "ffmpeg",
		FFprobePath:    "ffprobe",
	}

	// 创建媒体播放器实例
	mediaPlayer, err := NewMediaPlayer(cfg)
	if err != nil {
		log.Fatalf("Failed to create media player: %v", err)
	}

	// 示例：上传并处理视频
	if len(os.Args) > 1 {
		videoPath := os.Args[1]
		qualities := []VideoQuality{
			{Resolution: "1280x720", Bitrate: "2000k"},
			{Resolution: "854x480", Bitrate: "1000k"},
			{Resolution: "640x360", Bitrate: "500k"},
		}

		fmt.Printf("Processing video: %s\n", videoPath)
		if err := mediaPlayer.UploadAndProcessVideo(videoPath, "sample_video", qualities); err != nil {
			log.Printf("Error processing video: %v", err)
		} else {
			fmt.Printf("VOD URL: %s\n", mediaPlayer.GetVODURL("processed/sample_video/playlist.m3u8"))
		}
	}

	// 启动服务器
	fmt.Printf("Media server starting on port %s...\n", cfg.ServerPort)
	log.Fatal(mediaPlayer.Start())
}
