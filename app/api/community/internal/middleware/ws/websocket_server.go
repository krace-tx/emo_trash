package ws

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/krace-tx/emo_trash/app/api/community/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
	"net/http"
)

// WebSocketServer 结构体
type WebSocketServer struct {
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	manager *WebSocketManager
}

func NewWebSocketServer(svc *svc.ServiceContext) *WebSocketServer {
	return &WebSocketServer{
		ctx:     context.Background(),
		svcCtx:  svc,
		manager: NewWebSocketManager(),
	}
}

func (s *WebSocketServer) HandleConnection(w http.ResponseWriter, r *http.Request) {
	// 升级连接为 WebSocket
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		logx.Errorf("Upgrade connection failed: %v", err)
		return
	}

	// 从请求中获取用户 ID
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		logx.Errorf("User ID is empty")
		return
	}

	// 添加连接到管理器
	s.manager.addConnection(userID, conn)

	// 启动心跳检测
	threading.GoSafe(func() {
		s.manager.Heartbeat(s.ctx)
	})

	// 消息监听
	go func() {
		for {
			mt, message, err := conn.ReadMessage()
			if err != nil {
				logx.Errorf("Read message failed: %v", err)
				return
			}
			// 处理消息
			logx.Infof("Received message: %s", message)
			// 发送消息
			err = conn.WriteMessage(mt, message)
			if err != nil {
				logx.Errorf("Write message failed: %v", err)
				return
			}
		}
	}()
}
