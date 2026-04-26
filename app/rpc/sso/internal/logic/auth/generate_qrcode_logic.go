package authlogic

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/sso/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GenerateQrcodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGenerateQrcodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateQrcodeLogic {
	return &GenerateQrcodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 生成二维码
func (l *GenerateQrcodeLogic) GenerateQrcode(in *pb.GenerateQrcodeReq) (*pb.GenerateQrcodeResp, error) {
	qid := uuid.New().String()

	// 在 Redis 中初始化二维码状态为 WAITING
	// 有效期 2 分钟
	redisKey := fmt.Sprintf("qr:%s", qid)
	_ = l.svcCtx.Redis.Set(redisKey, "WAITING", 120*time.Second)

	l.Logger.Infof("生成登录二维码: qid=%s, device_id=%s", qid, in.DeviceId)

	// 实际项目中，这里返回一个能够生成二维码图片的 URL
	// 或者直接让前端拿 Qid 去生成
	return &pb.GenerateQrcodeResp{
		Qid:      qid,
		ImageUrl: fmt.Sprintf("https://emo-trash.app/qr/%s", qid),
	}, nil
}
