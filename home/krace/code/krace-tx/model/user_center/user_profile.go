package model

import "time"

type UserProfile struct {
	UserID        int64     `gorm:"primaryKey;column:user_id"        json:"user_id"`        // 用户唯一标识（关联UserAuth）
	Nickname      string    `gorm:"column:nickname;not null"         json:"nickname"`       // 用户昵称
	Username      string    `gorm:"column:username;unique;not null"  json:"username"`       // 唯一用户名（用于URL/提及，不可修改）
	Avatar        string    `gorm:"column:avatar"                   json:"avatar"`          // 头像URL
	CoverImage    string    `gorm:"column:cover_image"               json:"cover_image"`    // 封面图URL（个人主页顶部横幅）
	Bio           string    `gorm:"column:bio"                       json:"bio"`            // 个人简介(可选)
	Signature     string    `gorm:"column:signature"                 json:"signature"`      // 用户签名（短状态文案）
	Gender        int32     `gorm:"column:gender;default:null"       json:"gender"`         // 性别: 未知(null)/1-男/女-2（默认null更隐私）
	Birthdate     time.Time `gorm:"column:birthdate"                  json:"birthdate"`     // 生日时间戳（可选，用于年龄展示）
	Region        string    `gorm:"column:region"                    json:"region"`         // 地区(可选，如"北京")
	SocialLinks   string    `gorm:"column:social_links"              json:"social_links"`   // 社交链接(JSON格式，如{"weibo":"url","github":"url"})
	IsVerified    bool      `gorm:"column:is_verified;default:false" json:"is_verified"`    // 是否认证用户（如官方账号）
	VerifiedLabel string    `gorm:"column:verified_label"            json:"verified_label"` // 认证标签（如"官方认证"）
	ViewCount     int64     `gorm:"column:view_count;default:0"      json:"view_count"`     // 个人主页访问量
	CreatedAt     time.Time `gorm:"column:created_at"                json:"created_at"`     // 创建时间戳
	UpdatedAt     time.Time `gorm:"column:updated_at"                json:"updated_at"`     // 更新时间戳
}
