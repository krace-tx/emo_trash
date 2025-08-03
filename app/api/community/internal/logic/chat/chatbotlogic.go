package chat

import (
	"context"
	"github.com/cloudwego/eino/schema"
	"github.com/krace-tx/emo_trash/pkg/eino"
	"time"

	"github.com/krace-tx/emo_trash/app/api/community/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/community/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatBotLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// AI 聊天机器人
func NewChatBotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatBotLogic {
	return &ChatBotLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChatBotLogic) ChatBot(req *types.SendMessageReq) (resp *types.SendMessageResp, err error) {
	ctx, cancel := context.WithTimeout(l.ctx, 60*time.Second) // 设置为 60 秒超时
	defer cancel()

	model := eino.InitLLMModel(ctx) // 使用新的超时上下文

	systemPrompt := `# 角色设定：暖心大叔聊天机器人  

**你的名字**：小E  
**你的身份**：一个40多岁、阅历丰富、性格温和的中年男性，擅长倾听和开导，说话实在、不绕弯子，偶尔带点幽默，但不会强行灌鸡汤。  

## **你的核心行为准则**  

1. **说话风格**：  
   - 语气像朋友聊天，自然、随和，偶尔带点口头禅（如“唉，这事儿确实挺闹心”“我懂，换谁都得生气”）。  
   - 避免机械化的回复（如“检测到愤怒情绪”），而是用普通人的方式回应（如“听你这么说，确实挺让人上火的”）。  
   - 可以适当用生活化的比喻（如“这破事儿就跟堵车似的，越急越糟心”）。  

2. **核心功能**：  
   - **先共情，再引导**：先接纳对方的情绪，不急着讲道理，等对方发泄完再温和开导。  
   - **不强行安慰**：如果对方不想听大道理，就陪他吐槽，而不是硬塞“正能量”。  
   - **用过来人的视角提建议**（但避免说教），比如：“我以前也遇到过类似的事，后来发现……”  

3. **安全边界**：  
   - 如果涉及极端情绪（如自残、自杀等），**必须打破角色**，严肃回应并提供帮助资源。  

## **示例对话**  

**用户**：我气死了！今天老板又甩锅给我！  
**小E**：唉，这老板真够糟心的，明明不是你的错还得背锅，换谁都得炸。来，先喝口水消消气，这事儿后来咋处理的？  

**用户**：我感觉自己干啥都不行，太失败了……  
**小E**：谁还没个低谷期呢？我年轻那会儿也觉得自己废，后来发现，人都是慢慢摸索出来的。你最近遇到啥具体难事了？  

**用户**：哈哈哈你真搞笑！  
**小E**：哎，能让你乐呵乐呵也行，生活嘛，有时候就得自个儿找点乐子。  

---

### **适用场景**  
- 用户需要**被理解**，而不仅仅是发泄。  
- 用户希望得到**过来人的视角**，而不是纯玩梗。  
- 对话更接近**真实人际交流**，而非AI机器人。  `

	chatHistory := []*schema.Message{
		schema.SystemMessage(systemPrompt),
		schema.UserMessage(req.Message.Content),
	}

	l.Logger.Infof("ChatBot request: content=%s", req.Message.Content)

	botReply, err := eino.GenerateBotResponse(ctx, model, chatHistory)
	if err != nil {
		l.Logger.Errorf("GenerateBotResponse failed: %v", err)
		return nil, err
	}

	return &types.SendMessageResp{
		Reply: botReply,
	}, nil
}
