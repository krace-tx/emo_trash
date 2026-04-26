package postlogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/model"
	"github.com/krace-tx/emo_trash/app/rpc/post/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/post/pb"
	errx "github.com/krace-tx/emo_trash/pkg/err"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCommentLogic {
	return &DeleteCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除评论
func (l *DeleteCommentLogic) DeleteComment(in *pb.DeleteCommentReq) (*pb.CommonResp, error) {
	l.Logger.Infof("尝试删除评论: comment_id=%s, user_id=%s", in.Id, in.UserId)

	oid, err := primitive.ObjectIDFromHex(in.Id)
	if err != nil {
		l.Logger.Errorf("非法评论ID: %s", in.Id)
		return nil, errx.ErrSystemArgInvalid
	}

	commentColl := l.svcCtx.Mongo.Collection(model.CommentCollectionName)

	// 1. 查询评论确保它存在，并验证所有权
	var comment model.Comment
	err = commentColl.FindOne(l.ctx, bson.M{"_id": oid}).Decode(&comment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// 如果记录不存在，直接返回成功或自定义错误
			return nil, errx.New(errx.ErrSystemArgInvalid.Code, "回声不存在或已被删除")
		}
		l.Logger.Errorf("查询评论失败: %v, comment_id=%s", err, in.Id)
		return nil, errx.ErrDBQueryFailed
	}

	if comment.UserId != in.UserId {
		l.Logger.Errorf("越权删除评论: comment_id=%s, requester=%s, owner=%s", in.Id, in.UserId, comment.UserId)
		return nil, errx.New(errx.ErrSystemArgInvalid.Code, "无权删除该回声")
	}

	// 2. 删除评论
	_, err = commentColl.DeleteOne(l.ctx, bson.M{"_id": oid})
	if err != nil {
		l.Logger.Errorf("删除评论记录失败: %v, comment_id=%s", err, in.Id)
		return nil, errx.ErrDBDeleteFailed
	}

	// 3. 同步减少帖子的评论计数
	postOid, err := primitive.ObjectIDFromHex(comment.PostId)
	if err == nil {
		postColl := l.svcCtx.Mongo.Collection(model.PostCollectionName)
		_, _ = postColl.UpdateOne(l.ctx, bson.M{"_id": postOid}, bson.M{"$inc": bson.M{"comment_count": -1}})
	}

	l.Logger.Infof("评论删除成功: comment_id=%s, user_id=%s", in.Id, in.UserId)
	return &pb.CommonResp{
		Success: true,
		Message: "回声已删除",
	}, nil
}
