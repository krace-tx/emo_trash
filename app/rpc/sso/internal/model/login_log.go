package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// LoginLog 登录日志模型，记录用户登录行为相关信息
type LoginLog struct {
	ID         primitive.ObjectID     `json:"id" bson:"_id,omitempty"`                // 登录日志ID (MongoDB ObjectID)
	UserID     uint64                 `json:"user_id" bson:"user_id,omitempty"`       // 用户ID (MongoDB ObjectID)
	Account    string                 `json:"account" bson:"account"`                 // 登录账号 (用户名/手机号/邮箱等)
	Platform   string                 `json:"platform" bson:"platform"`               // 登录平台 (web/ios/android/小程序等)
	DeviceID   string                 `json:"device_id" bson:"device_id,omitempty"`   // 设备唯一标识
	IP         string                 `json:"ip" bson:"ip"`                           // 登录IP地址
	UserAgent  string                 `json:"user_agent" bson:"user_agent,omitempty"` // 用户代理信息 (浏览器/客户端版本等)
	LoginTime  time.Time              `json:"login_time" bson:"login_time"`           // 登录时间
	Status     string                 `json:"status" bson:"status"`                   // 登录状态 (success:成功/fail:失败)
	ResultCode int                    `json:"result_code" bson:"result_code"`         // 结果状态码 (错误码，成功时通常为0)
	ResultMsg  string                 `json:"result_msg" bson:"result_msg"`           // 结果描述信息 (成功/失败原因)
	TraceID    string                 `json:"trace_id" bson:"trace_id,omitempty"`     // 分布式追踪ID (用于链路追踪)
	Extra      map[string]interface{} `json:"extra" bson:"extra,omitempty"`           // 额外信息 (存储特殊场景下的扩展数据)
}

func (LoginLog) TableName() string {
	return "sso_login_log"
}
