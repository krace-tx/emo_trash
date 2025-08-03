package ws

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"log"
	"net/http"
	"sync"
	"time"
)

// WebSocketManager 管理连接、P2P、群聊等
type WebSocketManager struct {
	clients       map[string]*websocket.Conn // user_id -> conn
	groups        map[string]map[string]bool // group_id -> {user_id: true}
	offlineMsgs   map[string][]string        // user_id -> offline messages
	mu            sync.RWMutex               // 读写锁
	broadcast     chan string
	heartbeatTime time.Duration // 心跳检测时间间隔
}

// NewWebSocketManager 创建 WebSocketManager
func NewWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
		clients:       make(map[string]*websocket.Conn),
		groups:        make(map[string]map[string]bool),
		offlineMsgs:   make(map[string][]string),
		broadcast:     make(chan string),
		heartbeatTime: 10 * time.Second,
	}
}

// UpgradeConnection 升级 WebSocket 连接
func (wm *WebSocketManager) UpgradeConnection(userID string, w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return nil, err
	}

	if userID == "" {
		log.Println("User ID is required")
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return nil, fmt.Errorf("user ID is required")
	}

	// 添加用户连接
	wm.addConnection(userID, conn)

	return conn, nil
}

// addConnection 添加用户连接
func (wm *WebSocketManager) addConnection(userID string, conn *websocket.Conn) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	// 如果已有连接，关闭旧连接
	if oldConn, exists := wm.clients[userID]; exists {
		oldConn.Close()
	}
	wm.clients[userID] = conn
	log.Printf("User %s connected\n", userID)

	// 发送离线消息
	if messages, exists := wm.offlineMsgs[userID]; exists {
		for _, msg := range messages {
			conn.WriteMessage(websocket.TextMessage, []byte(msg))
		}
		delete(wm.offlineMsgs, userID)
	}
}

// CloseConnection 关闭用户连接
func (wm *WebSocketManager) CloseConnection(userID string) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	if conn, ok := wm.clients[userID]; ok {
		conn.Close()
		delete(wm.clients, userID)
		logx.Infof("Connection closed for user %s\n", userID)
	}
}

// GetConnection 根据 userID 获取连接
func (wm *WebSocketManager) GetConnection(userID string) *websocket.Conn {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	// 返回对应的连接
	if conn, exists := wm.clients[userID]; exists {
		return conn
	}

	// 用户连接不存在，返回 nil
	return nil
}

// IsConnectionActive 检查 WebSocket 连接是否活跃
func (wm *WebSocketManager) IsConnectionActive(userID string, conn *websocket.Conn) bool {
	defer func() {
		if r := recover(); r != nil {
			logx.Errorf("Recovered from panic in connection check for user %s: %v", userID, r)
			wm.CloseConnection(userID) // 确保异常连接被正确关闭
		}
	}()

	// 发送 Ping 消息测试连接
	if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
		logx.Errorf("Connection check failed: %v\n", err)
		return false
	}

	// 监听 Pong 消息或读超时检查
	conn.SetReadDeadline(time.Now().Add(wm.heartbeatTime)) // 设置超时时间
	_, _, err := conn.ReadMessage()
	if websocket.IsUnexpectedCloseError(err) {
		logx.Errorf("Unexpected close error: %v\n", err)
		logx.Infof("Connection for user %s is inactive. Closing it.", userID)
		wm.CloseConnection(userID)
		return false
	}
	return true
}

// Heartbeat 定期检测所有连接的活跃状态
func (wm *WebSocketManager) Heartbeat(ctx context.Context) {
	ticker := time.NewTicker(wm.heartbeatTime)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logx.Info("Heartbeat stopped.")
			return
		case <-ticker.C:
			wm.mu.RLock()
			for userID, conn := range wm.clients {
				go wm.IsConnectionActive(userID, conn)
			}
			wm.mu.RUnlock()
		}
	}
}
