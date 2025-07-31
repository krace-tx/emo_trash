package sociallogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/users/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/users/users"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReportUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReportUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReportUserLogic {
	return &ReportUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 举报用户
func (l *ReportUserLogic) ReportUser(in *users.ReportUserRequest) (*users.ReportUserResponse, error) {
	// todo: add your logic here and delete this line

	return &users.ReportUserResponse{}, nil
}
