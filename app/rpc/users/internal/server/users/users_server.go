// Code generated by goctl. DO NOT EDIT.
// goctl 1.8.5
// Source: users.proto

package server

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/users/internal/logic/users"
	"github.com/krace-tx/emo_trash/app/rpc/users/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/users/users"
)

type UsersServer struct {
	svcCtx *svc.ServiceContext
	users.UnimplementedUsersServer
}

func NewUsersServer(svcCtx *svc.ServiceContext) *UsersServer {
	return &UsersServer{
		svcCtx: svcCtx,
	}
}

// 获取用户的信息
func (s *UsersServer) GetUserInfo(ctx context.Context, in *users.GetUserInfoRequest) (*users.GetUserInfoResponse, error) {
	l := userslogic.NewGetUserInfoLogic(ctx, s.svcCtx)
	return l.GetUserInfo(in)
}

// 更新用户信息
func (s *UsersServer) EditUserInfo(ctx context.Context, in *users.EditUserInfoRequest) (*users.EditUserInfoResponse, error) {
	l := userslogic.NewEditUserInfoLogic(ctx, s.svcCtx)
	return l.EditUserInfo(in)
}

// 注册用户
func (s *UsersServer) RegisterUser(ctx context.Context, in *users.RegisterUserRequest) (*users.RegisterUserResponse, error) {
	l := userslogic.NewRegisterUserLogic(ctx, s.svcCtx)
	return l.RegisterUser(in)
}

// 编辑隐私权限
func (s *UsersServer) EditUserSetting(ctx context.Context, in *users.EditUserSettingRequest) (*users.EditUserSettingResponse, error) {
	l := userslogic.NewEditUserSettingLogic(ctx, s.svcCtx)
	return l.EditUserSetting(in)
}

// 获取隐私权限信息
func (s *UsersServer) GetUserSetting(ctx context.Context, in *users.GetUserSettingRequest) (*users.GetUserSettingResponse, error) {
	l := userslogic.NewGetUserSettingLogic(ctx, s.svcCtx)
	return l.GetUserSetting(in)
}

// 搜索用户
func (s *UsersServer) SearchUsers(ctx context.Context, in *users.SearchUsersRequest) (*users.SearchUsersResponse, error) {
	l := userslogic.NewSearchUsersLogic(ctx, s.svcCtx)
	return l.SearchUsers(in)
}
