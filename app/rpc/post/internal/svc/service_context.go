package svc

import (
	"context"
	"time"

	"github.com/krace-tx/emo_trash/app/model"
	"github.com/krace-tx/emo_trash/app/rpc/post/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ServiceContext struct {
	Config config.Config
	Mongo  *mongo.Database
}

func NewServiceContext(c config.Config) *ServiceContext {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(c.Mongo.Url))
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		panic(err)
	}

	db := client.Database(c.Mongo.Db)

	// 初始化索引
	_ = model.EnsurePostIndexes(context.Background(), db)
	_ = model.EnsureSocialIndexes(context.Background(), db)

	return &ServiceContext{
		Config: c,
		Mongo:  db,
	}
}
