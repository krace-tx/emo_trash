package authlogic

import (
	"context"
	"fmt"

	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/sso/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckQrcodeStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckQrcodeStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckQrcodeStatusLogic {
	return &CheckQrcodeStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 检查二维码状态
func (l *CheckQrcodeStatusLogic) CheckQrcodeStatus(in *pb.CheckQrcodeStatusReq) (*pb.CheckQrcodeStatusResp, error) {
	redisKey := fmt.Sprintf("qr:%s", in.Qid)

	// 从 Redis 获取状态
	var val string
	err := l.svcCtx.Redis.Get(redisKey, &val)
	if err != nil {
		// Redis 中不存在或已过期
		return &pb.CheckQrcodeStatusResp{
			Status: "EXPIRED",
		}, nil
	}

	// 简单的状态机解析 (实际中可能存储为 JSON 包含更多信息如 Token)
	// 这里假设 CONFIRMED 状态时，Redis 中存储的是 "CONFIRMED:access_token:refresh_token"
	status := val
	var accessToken, refreshToken string

	if len(val) > 9 && val[:9] == "CONFIRMED" {
		status = "CONFIRMED"
		// 解析 Token，这里为了简单演示，先写死或只返回状态
		// 在实际 ConfirmQrcodeLogin 中需要将 Token 写入此 key
	}

	return &pb.CheckQrcodeStatusResp{
		Status:       status,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
