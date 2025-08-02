package chatlogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/chat/chat"
	"github.com/krace-tx/emo_trash/app/rpc/chat/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReceiveMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReceiveMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReceiveMessageLogic {
	return &ReceiveMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 接收消息
func (l *ReceiveMessageLogic) ReceiveMessage(in *chat.ReceiveMessageRequest) (*chat.ReceiveMessageResponse, error) {
	// todo: add your logic here and delete this line

	return &chat.ReceiveMessageResponse{}, nil
}
