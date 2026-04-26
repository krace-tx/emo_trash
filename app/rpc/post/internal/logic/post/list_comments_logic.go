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

type ListCommentsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCommentsLogic {
	return &ListCommentsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 评论列表
func (l *ListCommentsLogic) ListComments(in *pb.ListCommentsReq) (*pb.ListCommentsResp, error) {
	l.Logger.Infof("获取回声列表: post_id=%s, cursor=%s", in.PostId, in.Cursor)

	filter := bson.M{"post_id": in.PostId}
	if in.Cursor != "" {
		// 这里简单处理，实际应该用 ObjectID 或者时间戳
		// filter["_id"] = bson.M{"$lt": cursorID}
	}

	opts := options.Find().
		SetSort(bson.M{"created_at": -1}).
		SetLimit(int64(in.PageSize + 1))

	commentColl := l.svcCtx.Mongo.Collection(model.CommentCollectionName)
	cursor, err := commentColl.Find(l.ctx, filter, opts)
	if err != nil {
		l.Logger.Errorf("查询回声列表失败: %v, post_id=%s", err, in.PostId)
		return nil, errx.ErrDBQueryFailed
	}
	defer cursor.Close(l.ctx)

	var comments []model.Comment
	if err := cursor.All(l.ctx, &comments); err != nil {
		return nil, errx.ErrDBQueryFailed
	}

	hasMore := false
	if len(comments) > int(in.PageSize) {
		hasMore = true
		comments = comments[:in.PageSize]
	}

	var nextCursor string
	if len(comments) > 0 {
		nextCursor = comments[len(comments)-1].ID.Hex()
	}

	pbList := make([]*pb.CommentInfo, 0, len(comments))
	for _, c := range comments {
		pbList = append(pbList, &pb.CommentInfo{
			Id:           c.ID.Hex(),
			PostId:       c.PostId,
			AuthorId:     c.UserId,
			AuthorName:   "匿名拾荒者", // TODO: 获取用户信息
			AuthorAvatar: "/static/logo.png",
			Content:      c.Content,
			CreatedAt:    c.CreatedAt.Unix(),
		})
	}

	return &pb.ListCommentsResp{
		List:       pbList,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}
