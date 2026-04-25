package postlogic

import (
	"context"
	"errors"

	"github.com/krace-tx/emo_trash/app/model"
	"github.com/krace-tx/emo_trash/app/rpc/post/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/post/pb"
	errx "github.com/krace-tx/emo_trash/pkg/err"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPostDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPostDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPostDetailLogic {
	return &GetPostDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPostDetailLogic) GetPostDetail(in *pb.GetPostDetailReq) (*pb.GetPostDetailResp, error) {
	oid, err := primitive.ObjectIDFromHex(in.Id)
	if err != nil {
		return nil, errx.ErrSystemArgInvalid
	}

	coll := l.svcCtx.Mongo.Collection(model.PostCollectionName)
	filter := bson.M{"_id": oid, "deleted_at": bson.M{"$exists": false}}

	var p model.Post
	if err := coll.FindOne(l.ctx, filter).Decode(&p); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errx.ErrPostNotFound
		}
		return nil, errx.ErrDBQueryFailed
	}

	return &pb.GetPostDetailResp{
		Post: &pb.PostInfo{
			Id:           p.ID.Hex(),
			AuthorId:     p.AuthorId,
			AuthorName:   p.AuthorName,
			AuthorAvatar: p.AuthorAvatar,
			Title:        p.Title,
			Content:      p.Content,
			Images:       p.Images,
			AiEvaluation: p.AiEvaluation,
			LikeCount:    p.LikeCount,
			CommentCount: p.CommentCount,
			StarCount:    p.StarCount,
			CreatedAt:    p.CreatedAt.Unix(),
		},
	}, nil
}
