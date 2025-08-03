package eino

import (
	"context"
	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino/schema"
	"github.com/ollama/ollama/api"
	"time"
)

// 初始化ChatMode
func Init(ctx context.Context, config *ollama.ChatModelConfig) *ollama.ChatModel {
	model, err := ollama.NewChatModel(ctx, config)

	if err != nil {
		panic(err)
	}
	return model
}

// 初始化LLM模型（复用测试代码中的配置）
func InitLLMModel(ctx context.Context) *ollama.ChatModel {
	keepAlive := 10 * time.Minute
	thinking := false
	model, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		BaseURL:   "http://localhost:11434",
		Timeout:   60 * time.Second,
		Model:     "qwen3:1.7b",
		KeepAlive: &keepAlive,
		Thinking:  &thinking,
		Options: &api.Options{
			Runner: api.Runner{
				NumCtx:    4096,
				NumGPU:    1,
				NumThread: 4,
			},
			Temperature: 0.7,
			TopP:        0.9,
			TopK:        40,
			NumPredict:  1000,
		},
	})
	if err != nil {
		panic(err)
	}
	return model
}

// 生成机器人回复
func GenerateBotResponse(ctx context.Context, model *ollama.ChatModel, chatHistory []*schema.Message) (string, error) {
	stream, err := model.Stream(ctx, chatHistory)
	if err != nil {
		return "", err
	}
	defer stream.Close()

	var reply string
	for {
		chunk, err := stream.Recv()
		if err != nil {
			break
		}
		reply += chunk.Content
	}

	return reply, nil
}
