package sensitive

import (
	"context"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"net/url"
)

type COSClient struct {
	client *cos.Client
}

// 初始化 COS 客户端
func NewCOSClient(bucketURL, ciURL, secretID, secretKey string) *COSClient {
	bu, err := url.Parse(bucketURL)
	if err != nil {
		logx.Errorf("invalid bucket URL: %v", err)
		return nil
	}
	cu, err := url.Parse(ciURL)
	if err != nil {
		logx.Errorf("invalid CI URL: %v", err)
		return nil
	}
	b := &cos.BaseURL{BucketURL: bu, CIURL: cu}

	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretID,
			SecretKey: secretKey,
		},
	})

	return &COSClient{client: client}
}

// 批量图片审核（支持多个 Base64 图片）
func (c *COSClient) BatchImageAuditing(base64Images []string) (*cos.BatchImageAuditingJobResult, error) {
	// 构造批量审核请求的 Input 列表
	detectType := "Porn,Terrorism,Politics,Ads"
	inputOptions := make([]cos.ImageAuditingInputOptions, len(base64Images))
	for i, base64Image := range base64Images {
		inputOptions[i] = cos.ImageAuditingInputOptions{
			DataId:  fmt.Sprintf("%d", i+1), // 生成唯一的 DataId
			Content: base64Image,            // 设置 Base64 图片内容
		}
	}

	// 构造请求选项
	opt := &cos.BatchImageAuditingOptions{
		Input: inputOptions,
		Conf: &cos.ImageAuditingJobConf{
			DetectType: detectType,
		},
	}

	// 发送批量审核请求
	res, _, err := c.client.CI.BatchImageAuditing(context.Background(), opt)
	if err != nil {
		return nil, fmt.Errorf("batch image auditing failed: %v", err)
	}
	return res, nil
}

// 判断图片审核是否通过
func IsImageApproved(result *cos.BatchImageAuditingJobResult) bool {
	if result == nil || len(result.JobsDetail) == 0 {
		return false
	}

	for _, job := range result.JobsDetail {
		if job.State != "Success" {
			fmt.Printf("Job state is not Success: %s\n", job.State)
			return false
		}

		if job.Label != "Normal" && job.Score >= 72 {
			fmt.Println(job.Label)
			fmt.Println(job.Score)
			return false
		}
	}

	return true
}
