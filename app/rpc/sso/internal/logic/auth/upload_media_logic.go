package authlogic

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/sso/pb"
	errx "github.com/krace-tx/emo_trash/pkg/err"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadMediaLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUploadMediaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadMediaLogic {
	return &UploadMediaLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 媒体上传
func (l *UploadMediaLogic) UploadMedia(in *pb.UploadMediaReq) (*pb.UploadMediaResp, error) {
	l.Logger.Infof("UploadMedia user_id: %s, usage: %s, filename: %s", in.UserId, in.Usage, in.Filename)

	// 1. 模拟 AI 审核内容 (待接入 AI 服务)
	// if !l.svcCtx.AiService.ReviewImage(in.Content) {
	//     return nil, errx.ErrSystemArgInvalid.WithMsg("图片内容审核未通过")
	// }

	// 2. 确定本地存储路径
	uploadDir := filepath.Join("uploads", in.Usage)
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			l.Logger.Errorf("创建上传目录失败: %v", err)
			return nil, errx.ErrSystemInternal
		}
	}

	// 3. 生成唯一文件名
	ext := filepath.Ext(in.Filename)
	newFilename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	filePath := filepath.Join(uploadDir, newFilename)

	// 4. 写入本地文件
	if err := os.WriteFile(filePath, in.Content, 0644); err != nil {
		l.Logger.Errorf("写入上传文件失败: %v, path=%s", err, filePath)
		return nil, errx.ErrSystemInternal
	}

	// 5. 返回访问 URL (假定网关配置了静态文件服务访问 /uploads)
	accessUrl := fmt.Sprintf("/uploads/%s/%s", in.Usage, newFilename)
	l.Logger.Infof("文件上传成功: user_id=%s, url=%s", in.UserId, accessUrl)

	return &pb.UploadMediaResp{
		Url: accessUrl,
	}, nil
}
