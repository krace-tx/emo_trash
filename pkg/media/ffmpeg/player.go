package ffmpeg

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/krace-tx/emo_trash/pkg/oss"
)

// MediaPlayer 媒体播放器
type MediaPlayer struct {
	minioClient *oss.MinioClient
	ffmpeg      *Processor
	server      *http.Server
	streams     map[string]*LiveStream
	mu          sync.RWMutex
}

// Config 播放器配置
type Config struct {
	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string
	MinioBucket    string
	MinioUseSSL    bool
	ServerPort     string
	FFmpegPath     string
	FFprobePath    string
}

// NewMediaPlayer 创建媒体播放器
func NewMediaPlayer(cfg Config) (*MediaPlayer, error) {
	// 初始化MinIO客户端
	minioCfg := oss.Config{
		Endpoint:  cfg.MinioEndpoint,
		AccessKey: cfg.MinioAccessKey,
		SecretKey: cfg.MinioSecretKey,
		Bucket:    cfg.MinioBucket,
		UseSSL:    cfg.MinioUseSSL,
	}

	minioClient, err := oss.NewMinioClient(minioCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %v", err)
	}

	// 初始化FFmpeg处理器
	ffmpegProcessor := NewProcessor(cfg.FFmpegPath, cfg.FFprobePath)

	player := &MediaPlayer{
		minioClient: minioClient,
		ffmpeg:      ffmpegProcessor,
		streams:     make(map[string]*LiveStream),
	}

	// 设置HTTP服务器
	mux := http.NewServeMux()
	mux.HandleFunc("/stream/", player.streamHandler)
	mux.HandleFunc("/vod/", player.vodHandler)
	mux.HandleFunc("/health", player.healthHandler)

	player.server = &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: mux,
	}

	return player, nil
}

// Start 启动媒体服务器
func (m *MediaPlayer) Start() error {
	log.Printf("Media server started at http://localhost:%s\n", strings.Split(m.server.Addr, ":")[1])
	return m.server.ListenAndServe()
}

// Stop 停止媒体服务器
func (m *MediaPlayer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return m.server.Shutdown(ctx)
}

// UploadAndProcessVideo 上传并处理视频
func (m *MediaPlayer) UploadAndProcessVideo(localPath, objectName string, qualities []VideoQuality) error {
	// 上传原始视频
	file, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	if err := m.minioClient.UploadFile("original/"+objectName, file, fileInfo.Size()); err != nil {
		return err
	}

	// 创建临时目录处理视频
	tempDir := filepath.Join(os.TempDir(), "media_processing", objectName)
	os.MkdirAll(tempDir, 0755)
	defer os.RemoveAll(tempDir)

	outputPlaylist := "playlist.m3u8"

	// 转码为HLS
	if err := m.ffmpeg.CreateHLSStream(localPath, tempDir, outputPlaylist, qualities); err != nil {
		return err
	}

	// 上传处理后的文件
	files, err := os.ReadDir(tempDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(tempDir, file.Name())
			uploadPath := filepath.Join("processed", objectName, file.Name())

			f, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer f.Close()

			fileInfo, _ := f.Stat()
			if err := m.minioClient.UploadFile(uploadPath, f, fileInfo.Size()); err != nil {
				return err
			}
		}
	}

	return nil
}

// streamHandler 处理流媒体请求
func (m *MediaPlayer) streamHandler(w http.ResponseWriter, r *http.Request) {
	m.handleCORS(w)

	if r.Method == http.MethodOptions {
		return
	}

	streamID := strings.TrimPrefix(r.URL.Path, "/stream/")
	if streamID == "" {
		http.Error(w, "Stream ID required", http.StatusBadRequest)
		return
	}

	m.mu.RLock()
	stream, exists := m.streams[streamID]
	m.mu.RUnlock()

	if !exists {
		http.Error(w, "Stream not found", http.StatusNotFound)
		return
	}

	if strings.HasSuffix(r.URL.Path, ".m3u8") {
		stream.servePlaylist(w, r)
	} else if strings.HasSuffix(r.URL.Path, ".ts") {
		stream.serveSegment(w, r)
	}
}

// vodHandler 处理点播请求
func (m *MediaPlayer) vodHandler(w http.ResponseWriter, r *http.Request) {
	m.handleCORS(w)

	if r.Method == http.MethodOptions {
		return
	}

	filePath := strings.TrimPrefix(r.URL.Path, "/vod/")
	if filePath == "" {
		http.Error(w, "File path required", http.StatusBadRequest)
		return
	}

	reader, err := m.minioClient.DownloadFile(filePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer reader.(io.Closer).Close()

	// 设置正确的Content-Type
	if strings.HasSuffix(filePath, ".m3u8") {
		w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
	} else if strings.HasSuffix(filePath, ".ts") {
		w.Header().Set("Content-Type", "video/mp2t")
	} else if strings.HasSuffix(filePath, ".mp4") {
		w.Header().Set("Content-Type", "video/mp4")
	}

	io.Copy(w, reader)
}

// healthHandler 健康检查
func (m *MediaPlayer) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// handleCORS 处理跨域请求
func (m *MediaPlayer) handleCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Max-Age", "86400")
}

// GetStreamURL 获取流媒体URL
func (m *MediaPlayer) GetStreamURL(streamID string) string {
	return fmt.Sprintf("http://localhost:%s/stream/%s/playlist.m3u8",
		strings.Split(m.server.Addr, ":")[1], streamID)
}

// GetVODURL 获取点播URL
func (m *MediaPlayer) GetVODURL(filePath string) string {
	return fmt.Sprintf("http://localhost:%s/vod/%s",
		strings.Split(m.server.Addr, ":")[1], filePath)
}
