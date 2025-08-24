package ffmpeg

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// VideoFormat 视频格式
type VideoFormat string

const (
	FormatHLS  VideoFormat = "hls"
	FormatMP4  VideoFormat = "mp4"
	FormatWebM VideoFormat = "webm"
	FormatDASH VideoFormat = "dash"
)

// VideoQuality 视频质量
type VideoQuality struct {
	Resolution string
	Bitrate    string
}

// Processor FFmpeg处理器
type Processor struct {
	ffmpegPath  string
	ffprobePath string
}

// NewProcessor 创建FFmpeg处理器
func NewProcessor(ffmpegPath, ffprobePath string) *Processor {
	if ffmpegPath == "" {
		ffmpegPath = "ffmpeg"
	}
	if ffprobePath == "" {
		ffprobePath = "ffprobe"
	}
	return &Processor{
		ffmpegPath:  ffmpegPath,
		ffprobePath: ffprobePath,
	}
}

// Transcode 转码视频
func (p *Processor) Transcode(inputPath, outputPath string, format VideoFormat, quality VideoQuality) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	args := []string{
		"-i", inputPath,
		"-y", // 覆盖输出文件
	}

	// 添加质量参数
	if quality.Resolution != "" {
		args = append(args, "-s", quality.Resolution)
	}
	if quality.Bitrate != "" {
		args = append(args, "-b:v", quality.Bitrate)
	}

	// 添加格式特定参数
	switch format {
	case FormatHLS:
		args = append(args,
			"-hls_time", "10",
			"-hls_list_size", "0",
			"-hls_segment_filename", filepath.Join(filepath.Dir(outputPath), "segment_%03d.ts"),
			"-f", "hls",
		)
	case FormatMP4:
		args = append(args,
			"-c:v", "libx264",
			"-preset", "medium",
			"-c:a", "aac",
			"-movflags", "+faststart",
			"-f", "mp4",
		)
	case FormatWebM:
		args = append(args,
			"-c:v", "libvpx-vp9",
			"-c:a", "libopus",
			"-f", "webm",
		)
	}

	args = append(args, outputPath)

	cmd := exec.CommandContext(ctx, p.ffmpegPath, args...)
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// CreateHLSStream 创建HLS流
func (p *Processor) CreateHLSStream(inputPath, outputDir, outputPlaylist string, qualities []VideoQuality) error {
	var commands []*exec.Cmd
	ctx := context.Background()

	// 为每个质量等级创建转码任务
	for i, quality := range qualities {
		playlistName := fmt.Sprintf("stream_%d.m3u8", i)
		segmentPattern := fmt.Sprintf("segment_%d_%%03d.ts", i)

		args := []string{
			"-i", inputPath,
			"-vf", fmt.Sprintf("scale=%s", quality.Resolution),
			"-c:v", "libx264",
			"-b:v", quality.Bitrate,
			"-c:a", "aac",
			"-f", "hls",
			"-hls_time", "10",
			"-hls_list_size", "0",
			"-hls_segment_filename", filepath.Join(outputDir, segmentPattern),
			filepath.Join(outputDir, playlistName),
		}

		cmd := exec.CommandContext(ctx, p.ffmpegPath, args...)
		commands = append(commands, cmd)
	}

	// 执行所有转码命令
	for _, cmd := range commands {
		if err := cmd.Start(); err != nil {
			return err
		}
	}

	// 等待所有命令完成
	for _, cmd := range commands {
		if err := cmd.Wait(); err != nil {
			return err
		}
	}

	// 创建主播放列表
	return p.createMasterPlaylist(outputDir, outputPlaylist, qualities)
}

func (p *Processor) createMasterPlaylist(outputDir, outputPlaylist string, qualities []VideoQuality) error {
	file, err := os.Create(filepath.Join(outputDir, outputPlaylist))
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.WriteString("#EXTM3U\n")

	for i, quality := range qualities {
		bandwidth := "2000000" // 默认带宽
		if quality.Bitrate != "" {
			bandwidth = strings.ReplaceAll(quality.Bitrate, "k", "000")
		}

		writer.WriteString(fmt.Sprintf("#EXT-X-STREAM-INF:BANDWIDTH=%s,RESOLUTION=%s\n", bandwidth, quality.Resolution))
		writer.WriteString(fmt.Sprintf("stream_%d.m3u8\n", i))
	}

	return writer.Flush()
}

// GetVideoInfo 获取视频信息
func (p *Processor) GetVideoInfo(inputPath string) (string, error) {
	cmd := exec.Command(p.ffprobePath,
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		inputPath,
	)

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// ExtractThumbnail 提取缩略图
func (p *Processor) ExtractThumbnail(inputPath, outputPath string, timeOffset string) error {
	args := []string{
		"-i", inputPath,
		"-ss", timeOffset,
		"-vframes", "1",
		"-q:v", "2",
		"-y",
		outputPath,
	}

	cmd := exec.Command(p.ffmpegPath, args...)
	return cmd.Run()
}
