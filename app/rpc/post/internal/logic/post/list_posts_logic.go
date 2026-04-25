package postlogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/model"
	"github.com/krace-tx/emo_trash/app/rpc/post/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/post/pb"
	errx "github.com/krace-tx/emo_trash/pkg/err"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPostsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListPostsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPostsLogic {
	return &ListPostsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListPostsLogic) ListPosts(in *pb.ListPostsReq) (*pb.ListPostsResp, error) {
	coll := l.svcCtx.Mongo.Collection(model.PostCollectionName)

	filter := bson.M{
		"status":     1,
		"deleted_at": bson.M{"$exists": false},
	}

	if in.Cursor != "" {
		lastId, err := primitive.ObjectIDFromHex(in.Cursor)
		if err == nil {
			filter["_id"] = bson.M{"$lt": lastId}
		}
	}

	limit := int64(in.PageSize)
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"_id": -1})
	findOptions.SetLimit(limit + 1)

	cursor, err := coll.Find(l.ctx, filter, findOptions)
	if err != nil {
		l.Logger.Errorf("查询帖子列表失败: %v", err)
		return nil, errx.ErrDBQueryFailed
	}
	defer cursor.Close(l.ctx)

	var list []model.Post
	if err := cursor.All(l.ctx, &list); err != nil {
		return nil, errx.ErrDBQueryFailed
	}

	hasMore := false
	nextCursor := ""
	if int64(len(list)) > limit {
		hasMore = true
		list = list[:limit]
		nextCursor = list[len(list)-1].ID.Hex()
	}

	var pbList []*pb.PostInfo
	for _, p := range list {
		pbList = append(pbList, &pb.PostInfo{
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
		})
	}

	return &pb.ListPostsResp{
		List:       pbList,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}
