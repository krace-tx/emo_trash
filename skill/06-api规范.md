# API 实现规范

## 适用范围

`app/api/gateway/internal/logic/**`、`types/**` 及所有面向 HTTP 的契约定义。

## API 逻辑职责

API logic 只做以下四件事，不得越界：

1. **接收并校验请求 DTO**：来自 `types` 定义的请求结构体
2. **映射请求到 RPC 入参**：构造 proto request 对象
3. **调用 RPC 并转换结果**：将 RPC 响应翻译为统一 `CommonResp`
4. **保持薄层**：禁止在 API 层嵌入已属于 RPC 的领域策略检查

## 响应契约规则

统一响应封包字段：

| 字段 | 说明 |
|---|---|
| `code` | 业务状态码 |
| `success` | 操作成功标志 |
| `data` | 业务数据负载 |
| `message` | 人可读的说明信息 |

- **失败**：`return types.Error(err), nil`
- **成功**：`return types.Success(data), nil`
- **特殊场景**：使用 `resp.go` 中的 builder helper，保持语义一致

## 请求定义规则

- 使用 goctl / go-zero 风格的 struct tag：`required`、`optional`、`len`、`regexp`、`options`。
- 字段名与 tag 保持稳定，不得随意重命名以保证向后兼容。
- 重要字段添加注释以保持 API 文档一致性。

## 日志规则

- API 层只记录调用结果与 trace context，不重复记录 RPC 层已详细记录的失败原因。
- 禁止记录凭证或 token 原始值。

## 错误处理规则

- RPC 详细失败原因通过 `pkg/err` 解析路径规范化，API 层不得擅自重写业务状态码或 message，除非接口契约有明确约定。