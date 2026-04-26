// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package sso

import (
	"context"

	"github.com/krace-tx/emo_trash/app/api/gateway/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/gateway/internal/types"
	"github.com/krace-tx/emo_trash/app/rpc/sso/client/auth"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserStatsLogic {
	return &GetUserStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserStatsLogic) GetUserStats(req *types.GetUserStatsReq) (resp *types.CommonResp, err error) {
	userId := l.ctx.Value("user_id").(string)

	data, err := l.svcCtx.Sso.GetUserStats(l.ctx, &auth.GetUserStatsReq{
		UserId: userId,
	})
	if err != nil {
		l.Logger.Errorf("获取用户统计失败: %v, user_id=%s", err, userId)
		return types.Error(err), nil
	}

	return types.Success(types.UserStats{
		PostCount: data.PostCount,
		LikeCount: data.ResonanceCount,
		JoinDays:  data.JoinDays,
	}), nil
}
