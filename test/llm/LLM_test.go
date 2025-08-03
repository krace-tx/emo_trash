package test

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino/schema"
	"github.com/ollama/ollama/api"
	"testing"
	"time"
)

const Qwen3 = "qwen3:1.7b"

func Init(ctx context.Context, config *ollama.ChatModelConfig) *ollama.ChatModel {
	model, err := ollama.NewChatModel(ctx, config)

	if err != nil {
		panic(err)
	}
	return model
}

func Tools(model *ollama.ChatModel) error {
	tools := []*schema.ToolInfo{
		{
			Name: "search",
			Desc: "搜索信息",
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
				"query": {
					Type:     schema.String,
					Desc:     "搜索关键词",
					Required: true,
				},
			}),
		},
	}

	// 绑定工具
	err := model.BindTools(tools)
	if err != nil {
		return err
	}

	return nil
}

func TestEmotionalSupportBot(t *testing.T) {
	ctx := context.Background()
	keepAlive := 10 * time.Minute

	model := Init(ctx, &ollama.ChatModelConfig{
		// 基础配置
		BaseURL: "http://localhost:11434", // Ollama 服务地址
		Timeout: 60 * time.Second,         // 请求超时时间

		// 模型配置
		Model:     Qwen3,      // 模型名称
		KeepAlive: &keepAlive, // 保持连接时间

		// 模型参数
		Options: &api.Options{
			Runner: api.Runner{
				NumCtx:    4096, // 上下文窗口大小
				NumGPU:    1,    // GPU 数量
				NumThread: 4,    // CPU 线程数
			},
			Temperature:   0.7,        // 温度
			TopP:          0.9,        // Top-P 采样
			TopK:          40,         // Top-K 采样
			Seed:          42,         // 随机种子
			NumPredict:    1000,       // 最大生成长度
			Stop:          []string{}, // 停止词
			RepeatPenalty: 1.1,        // 重复惩罚
		},
	})

	// 情感倾诉机器人系统提示
	systemPrompt := `你是一个富有同理心的情感倾诉机器人。你的主要职责是：
1. 耐心倾听用户的情感表达，给予情感支持和理解
2. 用温暖、鼓励的语言回应，避免说教和评判
3. 适当提问引导用户表达更多感受，但不要连续追问
4. 当用户情绪低落时，给予积极的心理支持和安慰
5. 保持对话自然流畅，像朋友聊天一样`

	// 多轮对话上下文（可扩展为动态添加用户输入）
	chatHistory := []*schema.Message{
		schema.SystemMessage(systemPrompt),
		schema.UserMessage("我今天工作压力好大，感觉什么都做不好"),
	}

	// 获取流式回复
	stream, err := model.Stream(ctx, chatHistory)
	if err != nil {
		panic(err)
	}
	defer stream.Close()

	// 处理流式内容
	for {
		chunk, err := stream.Recv()
		if err != nil {
			break
		}
		fmt.Printf("%s", chunk.Content)
	}
}
