package postlogic

import (
	"context"
	"time"

	"github.com/krace-tx/emo_trash/app/model"
	"github.com/krace-tx/emo_trash/app/rpc/post/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/post/pb"
	errx "github.com/krace-tx/emo_trash/pkg/err"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreatePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePostLogic {
	return &CreatePostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreatePostLogic) CreatePost(in *pb.CreatePostReq) (*pb.CreatePostResp, error) {
	now := time.Now()

	p := &model.Post{
		AuthorId:    in.AuthorId,
		Title:       in.Title,
		Content:     in.Content,
		Images:      in.Images,
		IsAnonymous: in.IsAnonymous,
		Status:      1, // 正常
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// 模拟 AI 评价注入
	p.AiEvaluation = "来自智能星工坊的评价：这是一份独特的情绪投递，我们已经收到了。"

	res, err := l.svcCtx.Mongo.Collection(model.PostCollectionName).InsertOne(l.ctx, p)
	if err != nil {
		l.Logger.Errorf("发布帖子失败: %v, author_id=%s", err, in.AuthorId)
		return nil, errx.ErrPostCreateFailed
	}

	oid := res.InsertedID.(primitive.ObjectID)
	l.Logger.Infof("用户发布帖子成功: user_id=%s, post_id=%s", in.AuthorId, oid.Hex())

	return &pb.CreatePostResp{
		Id: oid.Hex(),
	}, nil
}
