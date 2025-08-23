package model

import (
	"context"
	"strconv"

	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/svc"
	"github.com/krace-tx/emo_trash/pkg/db/rdb"
	errx "github.com/krace-tx/emo_trash/pkg/err"
	"github.com/zeromicro/go-zero/core/logx"
)

type UserModel struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserModel(ctx context.Context, svcCtx *svc.ServiceContext) *UserModel {
	return &UserModel{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CreateUser 创建用户数据
func (m *UserModel) CreateUser(userID uint64, mobile, password, salt string) (*UserAuth, *UserProfile, error) {
	var auth *UserAuth
	var profile *UserProfile

	transaction := func(engine rdb.EngineInterface[UserAuth]) error {
		auth = &UserAuth{
			UserID:   userID,
			Mobile:   mobile,
			Account:  mobile,
			Password: password,
			Salt:     salt,
			Status:   UserStatusNormal,
		}
		if err := engine.Create(m.ctx, auth); err != nil {
			return errx.ErrAuthCreateUserAuthFail
		}

		profileEngine := rdb.NewEngine[UserProfile](rdb.M)
		profile = &UserProfile{
			UserID:   userID,
			Nickname: "用户" + strconv.FormatUint(userID, 10)[8:],
			Avatar:   "",
		}
		if err := profileEngine.Create(m.ctx, profile); err != nil {
			return errx.ErrAuthCreateUserProfileFail
		}

		return nil
	}

	err := rdb.Transaction[UserAuth](m.ctx, rdb.M, transaction)
	if err != nil {
		return nil, nil, err
	}

	return auth, profile, nil
}
