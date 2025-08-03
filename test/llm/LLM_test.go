package test

import (
	"context"
	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/schema"
	"testing"
)

func Init(ctx context.Context, config *ark.ChatModelConfig) *ark.ChatModel {
	model, err := ark.NewChatModel(ctx, config)
	if err != nil {
		panic(err)
	}
	return model
}

func TestLLMModel(t *testing.T) {
	ctx := context.Background()
	config := &ark.ChatModelConfig{
		Model:  "llama2",
		APIKey: "123456",
	}
	chatModel := Init(ctx, config)

	input := []*schema.Message{
		schema.SystemMessage("你是一个专业的情感分析模型"),
		schema.UserMessage("我不喜欢这个东西"),
	}

	output, err := chatModel.Generate(ctx, input)
	if err != nil {
		panic(err)
	}
	t.Log(output.Content)

}
