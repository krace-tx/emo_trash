package sso

import (
	"context"

	"github.com/krace-tx/emo_trash/app/api/gateway/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/gateway/internal/types"
	"github.com/krace-tx/emo_trash/app/rpc/sso/client/auth"
	consts "github.com/krace-tx/emo_trash/pkg/constant"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoLogic {
	return &UpdateUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserInfoLogic) UpdateUserInfo(req *types.UpdateUserInfoReq) (resp *types.CommonResp, err error) {
	userId := l.ctx.Value(consts.UserId).(string)

	data, err := l.svcCtx.Sso.UpdateUserInfo(l.ctx, &auth.UpdateUserInfoReq{
		UserId:   userId,
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
		Bio:      req.Bio,
		Mood:     req.Mood,
	})
	if err != nil {
		l.Logger.Errorf("更新用户信息失败: %v, user_id=%s", err, userId)
		return types.Error(err), nil
	}

	return types.Success(data), nil
}
