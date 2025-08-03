package eino

import (
	"context"
	"github.com/cloudwego/eino-ext/components/model/ark"
)

// 初始化ChatMode
func InitChatModel(ctx context.Context) (*ark.ChatModel, error) {
	arkChatModel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		BaseURL: "http://localhost:11434",
		Model:   "llama2",
		APIKey:  "123456",
	})
	if err != nil {
		return nil, err
	}
	return arkChatModel, nil
}

func test(ctx context.Context, chatModel *ark.ChatModel) {

}
