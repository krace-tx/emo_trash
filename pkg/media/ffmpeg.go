package media

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"strings"
)

//go:embed video/src/*
var allFiles embed.FS

// CORS 处理函数：允许所有域访问
func handleCORS(w http.ResponseWriter) {
	// 允许所有来源进行跨域访问
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 允许的方法
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	// 允许的头部
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	// 缓存 CORS 预检请求的时间，避免每次请求都发送预检请求
	w.Header().Set("Access-Control-Max-Age", "86400")
}

// 处理 .m3u8 播放列表请求
func playlistHandler(w http.ResponseWriter, r *http.Request) {
	// 处理 CORS
	handleCORS(w)

	// 读取嵌入的 .m3u8 文件
	playlist, err := allFiles.ReadFile("video/src/playlist.m3u8")
	if err != nil {
		http.Error(w, "文件未找到", http.StatusNotFound)
		log.Println("Error reading playlist:", err)
		return
	}

	// 设置 Content-Type
	w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")

	// 返回嵌入的 playlist 文件内容
	_, err = w.Write(playlist)
	if err != nil {
		log.Println("Error sending playlist:", err)
	}
}

// 处理 .ts 视频片段请求
func tsHandler(w http.ResponseWriter, r *http.Request) {
	// 处理 CORS
	handleCORS(w)

	// 获取请求的 .ts 文件名
	fileName := strings.TrimPrefix(r.URL.Path, "/")

	// 读取嵌入的 .ts 文件
	tsFile, err := allFiles.ReadFile("video/src/" + fileName)
	if err != nil {
		http.Error(w, "文件未找到", http.StatusNotFound)
		log.Println("Error reading TS file:", err)
		return
	}

	// 设置 Content-Type
	w.Header().Set("Content-Type", "video/mp2t")

	// 返回嵌入的 .ts 文件内容
	_, err = w.Write(tsFile)
	if err != nil {
		log.Println("Error sending TS file:", err)
	}
}

// test
func main() {
	// 设置视频流的根目录
	http.HandleFunc("/playlist.m3u8", playlistHandler)
	http.HandleFunc("/", tsHandler) // 用于处理所有 .ts 文件请求

	// 启动 HTTP 服务器
	port := "8080"
	fmt.Printf("Server started at http://localhost:%s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Server failed: %v\n", err)
	}
}
