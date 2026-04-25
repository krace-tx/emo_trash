package model

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const PostCollectionName = "posts"

type Post struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"              json:"id"`
	AuthorId     string             `bson:"author_id"                  json:"author_id"`     // 作者ID
	AuthorName   string             `bson:"author_name"                json:"author_name"`   // 匿名昵称
	AuthorAvatar string             `bson:"author_avatar"              json:"author_avatar"` // 匿名头像
	Title        string             `bson:"title"                      json:"title"`         // 标题
	Content      string             `bson:"content"                    json:"content"`       // 内容
	Images       []string           `bson:"images"                     json:"images"`        // 图片列表
	Video        string             `bson:"video"                      json:"video"`         // 视频URL
	IsAnonymous  bool               `bson:"is_anonymous"               json:"is_anonymous"`  // 是否匿名
	AiEvaluation string             `bson:"ai_evaluation"              json:"ai_evaluation"` // AI评价
	LikeCount    int64              `bson:"like_count"                 json:"like_count"`    // 点赞数
	CommentCount int64              `bson:"comment_count"              json:"comment_count"` // 评论数
	StarCount    int64              `bson:"star_count"                 json:"star_count"`    // 收藏数
	Status       int8               `bson:"status"                     json:"status"`        // 状态: 1=正常, 0=审核中, -1=已删除
	CreatedAt    time.Time          `bson:"created_at"                 json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at"                 json:"updated_at"`
	DeletedAt    *time.Time         `bson:"deleted_at,omitempty"       json:"-"`
}

func (Post) TableName() string {
	return PostCollectionName
}

func EnsurePostIndexes(ctx context.Context, db *mongo.Database) error {
	coll := db.Collection(PostCollectionName)
	indexModels := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "author_id", Value: 1}, {Key: "_id", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}, {Key: "deleted_at", Value: 1}, {Key: "_id", Value: -1}},
		},
	}
	_, err := coll.Indexes().CreateMany(ctx, indexModels)
	return err
}
