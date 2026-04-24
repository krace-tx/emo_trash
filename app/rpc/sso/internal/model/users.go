package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserStatus int8

const (
	UserStatusNormal   UserStatus = 1
	UserStatusDisabled UserStatus = 0
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"              json:"id"`
	Email     string             `bson:"email"                      json:"email"`
	Salt      string             `bson:"salt"                       json:"-"`
	Password  string             `bson:"password"                   json:"-"` // bcrypt hash，不对外暴露
	Nickname  string             `bson:"nickname"                   json:"nickname"`
	Avatar    string             `bson:"avatar"                     json:"avatar"` // 头像URL
	Status    UserStatus         `bson:"status"                     json:"status"` // 1=正常 0=禁用
	CreatedAt time.Time          `bson:"created_at"                 json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"                 json:"updated_at"`
	DeletedAt *time.Time         `bson:"deleted_at,omitempty"       json:"-"` // 软删除，nil 表示未删除
}
