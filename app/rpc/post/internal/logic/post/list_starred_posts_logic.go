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

type ListStarredPostsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListStarredPostsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListStarredPostsLogic {
	return &ListStarredPostsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 我收藏的帖子列表
func (l *ListStarredPostsLogic) ListStarredPosts(in *pb.ListStarredPostsReq) (*pb.ListPostsResp, error) {
	l.Logger.Infof("获取收藏情绪列表: user_id=%s, cursor=%s", in.UserId, in.Cursor)

	// 1. 先从 Star 集合中获取收藏记录
	starColl := l.svcCtx.Mongo.Collection(model.StarCollectionName)
	filter := bson.M{"user_id": in.UserId}

	opts := options.Find().
		SetSort(bson.M{"created_at": -1}).
		SetLimit(int64(in.PageSize + 1))

	cursor, err := starColl.Find(l.ctx, filter, opts)
	if err != nil {
		l.Logger.Errorf("查询收藏记录失败: %v, user_id=%s", err, in.UserId)
		return nil, errx.ErrDBQueryFailed
	}
	defer cursor.Close(l.ctx)

	var stars []model.Star
	if err := cursor.All(l.ctx, &stars); err != nil {
		return nil, errx.ErrDBQueryFailed
	}

	hasMore := false
	if len(stars) > int(in.PageSize) {
		hasMore = true
		stars = stars[:in.PageSize]
	}

	if len(stars) == 0 {
		return &pb.ListPostsResp{}, nil
	}

	// 2. 批量获取帖子详情
	postIds := make([]primitive.ObjectID, 0, len(stars))
	for _, s := range stars {
		oid, _ := primitive.ObjectIDFromHex(s.PostId)
		postIds = append(postIds, oid)
	}

	postColl := l.svcCtx.Mongo.Collection(model.PostCollectionName)
	postCursor, err := postColl.Find(l.ctx, bson.M{
		"_id":        bson.M{"$in": postIds},
		"deleted_at": bson.M{"$exists": false},
	})
	if err != nil {
		l.Logger.Errorf("批量查询帖子失败: %v", err)
		return nil, errx.ErrDBQueryFailed
	}
	defer postCursor.Close(l.ctx)

	var posts []model.Post
	if err := postCursor.All(l.ctx, &posts); err != nil {
		return nil, errx.ErrDBQueryFailed
	}

	// 建立 ID 映射保持顺序
	postMap := make(map[string]model.Post)
	for _, p := range posts {
		postMap[p.ID.Hex()] = p
	}

	pbList := make([]*pb.PostInfo, 0, len(stars))
	for _, s := range stars {
		if p, ok := postMap[s.PostId]; ok {
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
				IsStarred:    true,
				CreatedAt:    p.CreatedAt.Unix(),
			})
		}
	}

	var nextCursor string
	if len(stars) > 0 {
		nextCursor = stars[len(stars)-1].ID.Hex()
	}

	return &pb.ListPostsResp{
		List:       pbList,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}
