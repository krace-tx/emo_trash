package no_sql

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConf struct {
	URI      string // MongoDB连接URI (必填)
	Database string // 要使用的数据库名称 (必填)
	Timeout  int    `json:",optional"` // 连接超时时间(秒) (可选，默认5)
	MaxPool  uint64 `json:",optional"` // 最大连接池大小 (可选，默认100)
	MinPool  uint64 `json:",optional"` // 最小连接池大小 (可选，默认10)
}

// InitMongo 根据配置初始化MongoDB客户端和数据库实例
func InitMongo(conf MongoConf) (*mongo.Database, error) {
	// 创建客户端选项
	clientOpts := options.Client().ApplyURI(conf.URI)

	// 设置连接池参数
	if conf.MaxPool > 0 {
		clientOpts.SetMaxPoolSize(conf.MaxPool)
	} else {
		clientOpts.SetMaxPoolSize(100) // 默认最大连接池
	}
	if conf.MinPool > 0 {
		clientOpts.SetMinPoolSize(conf.MinPool)
	} else {
		clientOpts.SetMinPoolSize(10) // 默认最小连接池
	}

	// 设置连接超时
	timeout := time.Duration(conf.Timeout) * time.Second
	if timeout <= 0 {
		timeout = 5 * time.Second // 默认超时5秒
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// 建立连接
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}

	// 验证连接
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	// 返回客户端和指定数据库实例
	return client.Database(conf.Database), nil
}

func MustInitMongo(conf MongoConf) *mongo.Database {
	db, err := InitMongo(conf)
	if err != nil {
		panic(err)
	}
	return db
}
