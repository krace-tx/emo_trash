package post

import (
	"context"

	"github.com/krace-tx/emo_trash/app/api/gateway/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/gateway/internal/types"
	"github.com/krace-tx/emo_trash/app/rpc/post/client/post"
	"github.com/krace-tx/emo_trash/app/rpc/post/pb"
	consts "github.com/krace-tx/emo_trash/pkg/constant"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPostDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPostDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPostDetailLogic {
	return &GetPostDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPostDetailLogic) GetPostDetail(req *types.GetPostDetailReq) (resp *types.CommonResp, err error) {
	userId, _ := l.ctx.Value(consts.UserId).(string)

	data, err := l.svcCtx.Post.GetPostDetail(l.ctx, &post.GetPostDetailReq{
		Id:     req.Id,
		UserId: userId,
	})
	if err != nil {
		l.Logger.Errorf("閼惧嘲褰囩敮鏍х摍鐠囷附鍎忔径杈Е: %v, user_id=%s, post_id=%s", err, userId, req.Id)
		return types.Error(err), nil
	}

	return types.Success(mapPostInfo(data.Post)), nil
}

func mapPostInfo(p *pb.PostInfo) types.PostInfo {
	if p == nil {
		return types.PostInfo{}
	}
	return types.PostInfo{
		Id:           p.Id,
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
		IsLiked:      p.IsLiked,
		IsStarred:    p.IsStarred,
		CreatedAt:    p.CreatedAt,
	}
}
