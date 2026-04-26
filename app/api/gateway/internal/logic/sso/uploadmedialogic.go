// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package sso

import (
	"context"

	"github.com/krace-tx/emo_trash/app/api/gateway/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadMediaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadMediaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadMediaLogic {
	return &UploadMediaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadMediaLogic) UploadMedia(req *types.UploadMediaReq) (resp *types.CommonResp, err error) {
	// 1. 从 context 获取 r (需要 go-zero 的 handler 传递，或者直接在 logic 处理)
	// 在 goctl 生成的逻辑中，通常建议在 handler 层解析文件，这里为了简单演示调用 RPC
	// 注意：API 定义中 UploadMediaReq 只有 Usage，文件数据需要从 HTTP request 获取

	// 这里模拟调用 RPC，实际开发中需要从 r.FormFile 获取内容
	// 为了对齐需求文档：支持本地存储 + AI 审核

	userId := l.ctx.Value("user_id").(string)

	// 示例：此处假设我们已经通过某种方式获取了文件内容 content 和 filename
	// rpcResp, err := l.svcCtx.SsoRpc.UploadMedia(l.ctx, &pb.UploadMediaReq{
	// 	UserId: userId,
	// 	Usage:  req.Usage,
	// 	Content: content,
	// 	Filename: filename,
	// })

	l.Logger.Infof("User %s uploading media for %s", userId, req.Usage)

	return types.Success(types.UploadMediaResp{
		Url: "/uploads/placeholder.png", // 实际应返回 RPC 返回的路径
	}), nil
}
