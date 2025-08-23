package model

import (
	"context"
	"time"

	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/svc"
	"github.com/krace-tx/emo_trash/pkg/db/rdb"
	errx "github.com/krace-tx/emo_trash/pkg/err"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

var (
	UserStatusNormal   = int32(0)
	UserStatusDisabled = int32(1)
)

// 用户认证信息
type UserAuth struct {
	UserID      uint64     `gorm:"primaryKey;column:user_id" json:"user_id"`           // 用户唯一标识（关联UserProfile）
	Mobile      string     `gorm:"column:mobile"             json:"mobile"`            // 手机号(可选)
	Email       string     `gorm:"column:email"              json:"email"`             // 邮箱(可选)
	Account     string     `gorm:"column:account"            json:"account"`           // 账号(可选)
	Password    string     `gorm:"column:password;not null"  json:"password"`          // 密码(加密存储)
	Salt        string     `gorm:"column:salt;not null"      json:"salt"`              // 密码盐值
	Platform    string     `gorm:"column:platform"           json:"platform"`          // 第三方平台: wechat/qq(可选)
	OpenID      string     `gorm:"column:open_id"            json:"open_id"`           // 第三方平台OpenID(可选)
	Status      int32      `gorm:"column:status;default:0"   json:"status"`            // 用户状态: 0-正常, 1-禁用
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"` // 创建时间戳
	UpdatedAt   time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"` // 更新时间戳
	LastLoginAt *time.Time `gorm:"column:last_login_at;null" json:"last_login_at"`     // 最后登录时间戳 (允许NULL)
	LastLoginIP string     `gorm:"column:last_login_ip"      json:"last_login_ip"`     // 最后登录IP
}

func (UserAuth) TableName() string {
	return "sso_user_auth"
}

type UserProfile struct {
	UserID    uint64    `gorm:"primaryKey;column:user_id" json:"user_id"`   // 用户唯一标识（关联UserAuth）
	Nickname  string    `gorm:"column:nickname;not null"  json:"nickname"`  // 用户昵称
	Avatar    string    `gorm:"column:avatar"             json:"avatar"`    // 头像URL
	Bio       string    `gorm:"column:bio"                json:"bio"`       // 个人简介(可选)
	Gender    int32     `gorm:"column:gender;default:0"   json:"gender"`    // 性别: 0-未知, 1-男, 2-女
	Region    string    `gorm:"column:region"             json:"region"`    // 地区(可选，如"北京")
	CreatedAt time.Time `gorm:"column:created_at"        json:"created_at"` // 创建时间戳
	UpdatedAt time.Time `gorm:"column:updated_at"        json:"updated_at"` // 更新时间戳
}

func (UserProfile) TableName() string {
	return "sso_user_profile"
}

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
func (m *UserModel) CreateUser(account, mobile, password, salt string) (*UserAuth, *UserProfile, error) {
	var auth *UserAuth
	var profile *UserProfile

	transaction := func(tx *gorm.DB) error {
		auth = &UserAuth{
			UserID:   uint64(m.svcCtx.Snowflake.Generate()),
			Mobile:   mobile,
			Account:  account,
			Password: password,
			Salt:     salt,
			Status:   UserStatusNormal,
		}

		if err := rdb.NewEngine[UserAuth](tx).Create(m.ctx, auth); err != nil {
			return errx.ErrAuthCreateUserAuthFail
		}

		profile = &UserProfile{
			UserID:   auth.UserID,
			Nickname: account,
			Avatar:   account,
			Bio:      "这个人很懒，什么都没有写",
			Region:   "未知",
		}

		if err := rdb.NewEngine[UserProfile](tx).Create(m.ctx, profile); err != nil {
			return errx.ErrAuthCreateUserProfileFail
		}

		return nil
	}

	if err := rdb.Transaction(m.ctx, rdb.M, transaction); err != nil {
		return nil, nil, err
	}

	return auth, profile, nil
}
