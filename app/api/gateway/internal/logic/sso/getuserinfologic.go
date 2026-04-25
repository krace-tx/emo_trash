package sso

import (
	"context"

	"github.com/krace-tx/emo_trash/app/api/gateway/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/gateway/internal/types"
	"github.com/krace-tx/emo_trash/app/rpc/sso/client/auth"
	consts "github.com/krace-tx/emo_trash/pkg/constant"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserInfoLogic) GetUserInfo(req *types.GetUserInfoReq) (resp *types.CommonResp, err error) {
	userId := l.ctx.Value(consts.UserId).(string)

	data, err := l.svcCtx.Sso.GetUserInfo(l.ctx, &auth.GetUserInfoReq{
		UserId: userId,
	})
	if err != nil {
		l.Logger.Errorf("获取用户信息失败: %v, user_id=%s", err, userId)
		return types.Error(err), nil
	}

	return types.Success(&types.UserInfo{
		UserId:     data.User.UserId,
		Email:      data.User.Email,
		Nickname:   data.User.Nickname,
		Avatar:     data.User.Avatar,
		CreateTime: data.User.CreateTime,
		Bio:        data.User.Bio,
		Mood:       data.User.Mood,
	}), nil
}
