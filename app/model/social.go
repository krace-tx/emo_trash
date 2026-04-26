package model

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	LikeCollectionName    = "post_likes"
	StarCollectionName    = "post_stars"
	CommentCollectionName = "post_comments"
)

type Like struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PostId    string             `bson:"post_id"     json:"post_id"`
	UserId    string             `bson:"user_id"     json:"user_id"`
	CreatedAt time.Time          `bson:"created_at"  json:"created_at"`
}

type Star struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PostId    string             `bson:"post_id"     json:"post_id"`
	UserId    string             `bson:"user_id"     json:"user_id"`
	CreatedAt time.Time          `bson:"created_at"  json:"created_at"`
}

type Comment struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PostId    string             `bson:"post_id"     json:"post_id"`
	UserId    string             `bson:"user_id"     json:"user_id"`
	Content   string             `bson:"content"     json:"content"`
	CreatedAt time.Time          `bson:"created_at"  json:"created_at"`
}

func EnsureSocialIndexes(ctx context.Context, db *mongo.Database) error {
	likeColl := db.Collection(LikeCollectionName)
	starColl := db.Collection(StarCollectionName)
	commentColl := db.Collection(CommentCollectionName)

	unique := true
	likeIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "post_id", Value: 1}, {Key: "user_id", Value: 1}},
		Options: &options.IndexOptions{Unique: &unique},
	}
	starIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "post_id", Value: 1}, {Key: "user_id", Value: 1}},
		Options: &options.IndexOptions{Unique: &unique},
	}
	commentIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "post_id", Value: 1}, {Key: "created_at", Value: -1}},
	}

	if _, err := likeColl.Indexes().CreateOne(ctx, likeIndex); err != nil {
		return err
	}
	if _, err := starColl.Indexes().CreateOne(ctx, starIndex); err != nil {
		return err
	}
	if _, err := commentColl.Indexes().CreateOne(ctx, commentIndex); err != nil {
		return err
	}
	return nil
}
