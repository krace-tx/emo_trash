package model

import "time" // 新增时间包导入

var (
	UserStatusNormal   = int32(0)
	UserStatusDisabled = int32(1)
)

type UserAuth struct {
	UserID      uint64    `gorm:"primaryKey;column:user_id" json:"user_id"`              // 用户唯一标识（关联UserProfile）
	Mobile      string    `gorm:"column:mobile"             json:"mobile"`               // 手机号(可选)
	Email       string    `gorm:"column:email"              json:"email"`                // 邮箱(可选)
	Account     string    `gorm:"column:account"            json:"account"`              // 账号(可选)
	Password    string    `gorm:"column:password;not null"  json:"password"`             // 密码(加密存储)
	Salt        string    `gorm:"column:salt;not null"      json:"salt"`                 // 密码盐值
	Platform    string    `gorm:"column:platform"           json:"platform"`             // 第三方平台: wechat/qq(可选)
	OpenID      string    `gorm:"column:open_id"            json:"open_id"`              // 第三方平台OpenID(可选)
	Status      int32     `gorm:"column:status;default:0"   json:"status"`               // 用户状态: 0-正常, 1-禁用
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`    // 创建时间戳
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`    // 更新时间戳
	LastLoginAt time.Time `gorm:"column:last_login_at"             json:"last_login_at"` // 最后登录时间戳
	LastLoginIP string    `gorm:"column:last_login_ip"      json:"last_login_ip"`        // 最后登录IP
}
