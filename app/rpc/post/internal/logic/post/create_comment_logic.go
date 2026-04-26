package postlogic

import (
	"context"
	"time"

	"github.com/krace-tx/emo_trash/app/model"
	"github.com/krace-tx/emo_trash/app/rpc/post/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/post/pb"
	errx "github.com/krace-tx/emo_trash/pkg/err"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCommentLogic {
	return &CreateCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CreateComment 创建回声 (评论)
func (l *CreateCommentLogic) CreateComment(in *pb.CreateCommentReq) (*pb.CommonResp, error) {
	l.Logger.Infof("发表回声: post_id=%s, user_id=%s", in.PostId, in.UserId)

	pid, err := primitive.ObjectIDFromHex(in.PostId)
	if err != nil {
		l.Logger.Errorf("非法帖子ID: %s", in.PostId)
		return nil, errx.ErrSystemArgInvalid
	}

	// 1. 验证帖子是否存在
	postColl := l.svcCtx.Mongo.Collection(model.PostCollectionName)
	var post model.Post
	if err := postColl.FindOne(l.ctx, bson.M{
		"_id":        pid,
		"deleted_at": bson.M{"$exists": false},
	}).Decode(&post); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errx.ErrPostNotFound
		}
		l.Logger.Errorf("查询帖子失败: %v, post_id=%s", err, in.PostId)
		return nil, errx.ErrDBQueryFailed
	}

	// 2. 插入评论记录
	commentColl := l.svcCtx.Mongo.Collection(model.CommentCollectionName)
	_, err = commentColl.InsertOne(l.ctx, model.Comment{
		ID:        primitive.NewObjectID(),
		PostId:    in.PostId,
		UserId:    in.UserId,
		Content:   in.Content,
		CreatedAt: time.Now(),
	})
	if err != nil {
		l.Logger.Errorf("插入回声记录失败: %v, post_id=%s, user_id=%s", err, in.PostId, in.UserId)
		return nil, errx.ErrDBInsertFailed
	}

	// 3. 更新帖子的评论计数
	_, _ = postColl.UpdateOne(l.ctx, bson.M{"_id": pid}, bson.M{"$inc": bson.M{"comment_count": 1}})

	l.Logger.Infof("回声发表成功: post_id=%s, user_id=%s", in.PostId, in.UserId)
	return &pb.CommonResp{
		Success: true,
		Message: "回声已传达",
	}, nil
}
