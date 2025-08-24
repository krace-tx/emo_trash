package ffmpeg

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// LiveStream 直播流
type LiveStream struct {
	ID        string
	outputDir string
	playlist  string
	segments  []string
	mu        sync.RWMutex
	isLive    bool
	stopChan  chan struct{}
	ffmpeg    *Processor
}

// NewLiveStream 创建新的直播流
func (m *MediaPlayer) NewLiveStream(streamID string, inputPath string, qualities []VideoQuality) (*LiveStream, error) {
	// 创建临时输出目录
	outputDir := filepath.Join(os.TempDir(), "live_streams", streamID)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, err
	}

	stream := &LiveStream{
		ID:        streamID,
		outputDir: outputDir,
		playlist:  "playlist.m3u8",
		isLive:    true,
		stopChan:  make(chan struct{}),
		ffmpeg:    m.ffmpeg,
	}

	// 启动转码过程
	go stream.startTranscoding(inputPath, qualities)

	m.mu.Lock()
	m.streams[streamID] = stream
	m.mu.Unlock()

	return stream, nil
}

// startTranscoding 开始转码
func (ls *LiveStream) startTranscoding(inputPath string, qualities []VideoQuality) {
	// 使用FFmpeg创建HLS直播流
	//outputPlaylist := filepath.Join(ls.outputDir, ls.playlist)

	err := ls.ffmpeg.CreateHLSStream(inputPath, ls.outputDir, ls.playlist, qualities)
	if err != nil {
		log.Printf("Error starting live stream %s: %v", ls.ID, err)
		ls.isLive = false
		return
	}

	// 监控并清理旧的分段
	go ls.cleanupOldSegments()
}

// cleanupOldSegments 清理旧的分段
func (ls *LiveStream) cleanupOldSegments() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ls.stopChan:
			return
		case <-ticker.C:
			ls.mu.Lock()
			// 保留最近10个分段
			if len(ls.segments) > 10 {
				segmentsToRemove := ls.segments[:len(ls.segments)-10]
				for _, segment := range segmentsToRemove {
					os.Remove(filepath.Join(ls.outputDir, segment))
				}
				ls.segments = ls.segments[len(segmentsToRemove):]
			}
			ls.mu.Unlock()
		}
	}
}

// servePlaylist 服务播放列表
func (ls *LiveStream) servePlaylist(w http.ResponseWriter, r *http.Request) {
	if !ls.isLive {
		http.Error(w, "Stream not live", http.StatusNotFound)
		return
	}

	playlistPath := filepath.Join(ls.outputDir, ls.playlist)
	http.ServeFile(w, r, playlistPath)
}

// serveSegment 服务分段
func (ls *LiveStream) serveSegment(w http.ResponseWriter, r *http.Request) {
	if !ls.isLive {
		http.Error(w, "Stream not live", http.StatusNotFound)
		return
	}

	segmentName := filepath.Base(r.URL.Path)
	segmentPath := filepath.Join(ls.outputDir, segmentName)
	http.ServeFile(w, r, segmentPath)
}

// Stop 停止直播流
func (ls *LiveStream) Stop() {
	ls.isLive = false
	close(ls.stopChan)

	// 清理文件
	os.RemoveAll(ls.outputDir)
}

// IsLive 检查流是否正在直播
func (ls *LiveStream) IsLive() bool {
	return ls.isLive
}

// GetStats 获取流统计信息
func (ls *LiveStream) GetStats() map[string]interface{} {
	ls.mu.RLock()
	defer ls.mu.RUnlock()

	return map[string]interface{}{
		"id":           ls.ID,
		"is_live":      ls.isLive,
		"segments":     len(ls.segments),
		"output_dir":   ls.outputDir,
		"last_updated": time.Now(),
	}
}
