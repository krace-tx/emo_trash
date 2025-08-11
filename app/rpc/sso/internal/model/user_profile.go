package model

import "time"

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
