package postlogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/model"
	"github.com/krace-tx/emo_trash/app/rpc/post/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/post/pb"
	errx "github.com/krace-tx/emo_trash/pkg/err"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMyPostsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListMyPostsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMyPostsLogic {
	return &ListMyPostsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 我发布的帖子列表
func (l *ListMyPostsLogic) ListMyPosts(in *pb.ListMyPostsReq) (*pb.ListPostsResp, error) {
	l.Logger.Infof("获取我的情绪列表: user_id=%s, cursor=%s", in.UserId, in.Cursor)

	filter := bson.M{
		"author_id":  in.UserId,
		"deleted_at": bson.M{"$exists": false},
	}
	// TODO: Handle cursor

	opts := options.Find().
		SetSort(bson.M{"created_at": -1}).
		SetLimit(int64(in.PageSize + 1))

	postColl := l.svcCtx.Mongo.Collection(model.PostCollectionName)
	cursor, err := postColl.Find(l.ctx, filter, opts)
	if err != nil {
		l.Logger.Errorf("查询我的情绪列表失败: %v, user_id=%s", err, in.UserId)
		return nil, errx.ErrDBQueryFailed
	}
	defer cursor.Close(l.ctx)

	var posts []model.Post
	if err := cursor.All(l.ctx, &posts); err != nil {
		return nil, errx.ErrDBQueryFailed
	}

	hasMore := false
	if len(posts) > int(in.PageSize) {
		hasMore = true
		posts = posts[:in.PageSize]
	}

	var nextCursor string
	if len(posts) > 0 {
		nextCursor = posts[len(posts)-1].ID.Hex()
	}

	pbList := make([]*pb.PostInfo, 0, len(posts))
	for _, p := range posts {
		pbList = append(pbList, &pb.PostInfo{
			Id:           p.ID.Hex(),
			AuthorId:     p.AuthorId,
			AuthorName:   p.AuthorName,
			AuthorAvatar: p.AuthorAvatar,
			Title:        p.Title,
			Content:      p.Content,
			Images:       p.Images,
			Video:        p.Video,
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
