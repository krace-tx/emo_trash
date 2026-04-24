# emo_trash

`emo_trash` 是一个基于 Go + go-zero 的微服务项目。  
当前仓库已完成并对齐的核心链路是 **SSO 认证体系（邮箱验证码、注册、登录、密码、Token）**。

## 当前实现范围

- 网关 API：`app/api/gateway`
- SSO RPC：`app/rpc/sso`
- 协议定义：
  - `cmd/api/gateway.api`
  - `cmd/api/sso.api`
  - `cmd/rpc/sso.proto`

当前 SSO 对外能力（HTTP）：

- `POST /sso/v1/auth/email/code`
- `POST /sso/v1/auth/register`
- `POST /sso/v1/auth/login`
- `POST /sso/v1/auth/password/reset`
- `POST /sso/v1/auth/password/change`
- `POST /sso/v1/auth/token/refresh`
- `POST /sso/v1/auth/token/verify`
- `POST /sso/v1/auth/logout`

## 技术栈

- Go: `1.24.5`
- 框架: `go-zero`
- 通信: `HTTP + gRPC`
- 存储: `MySQL + MongoDB + Redis`
- 认证: `JWT`
- 校验: `protoc-gen-validate`

## 项目结构（核心）

```text
.
├── app/
│   ├── api/gateway/                 # HTTP 网关
│   └── rpc/sso/                     # SSO RPC 服务
├── cmd/
│   ├── api/gateway.api              # 网关路由定义
│   ├── api/sso.api                  # API DTO 定义
│   └── rpc/sso.proto                # RPC 协议定义
├── pkg/
│   ├── auth/                        # JWT / 密码工具
│   ├── datastore/                   # mysql / mongo / redis 封装
│   ├── email/                       # 邮件发送
│   ├── err/                         # 统一错误码
│   └── ...
├── deploy/                          # docker/离线部署相关
├── third_party/                     # proto 三方依赖
├── Makefile
└── README.md
```

## 配置说明

本项目已在 `.gitignore` 中忽略 `**/etc/*.yaml` 与 `**/etc/*.yml`。  
请在本地自行创建配置文件，不要提交真实凭据。

### 1) SSO RPC 配置（示例）

运行 `app/rpc/sso/sso.go` 时，默认读取 `etc/sso.yaml`。

最小示例（请按环境替换）：

```yaml
Service:
  Name: sso.rpc
  Mode: dev

Zrpc:
  Name: sso.rpc
  ListenOn: 0.0.0.0:8100
  Etcd:
    Hosts:
      - 127.0.0.1:2379
    Key: sso.rpc

Mysql:
  DSN: root:password@tcp(127.0.0.1:3306)/emo_trash?charset=utf8mb4&parseTime=True&loc=Local

Mongo:
  URI: "mongodb://root:password@127.0.0.1:27017/emo_trash?authSource=admin"
  Database: "emo_trash"

Redis:
  Addr: 127.0.0.1:6379
  Password: ""
  DB: 0

JWT:
  AccessSecret: "replace_access_secret"
  AccessExpire: 86400
  RefreshSecret: "replace_refresh_secret"
  RefreshExpire: 604800

Email:
  Id: "your_email@qq.com"
  Auth: "your_smtp_auth_code"
  Host: "smtp.qq.com"
  Port: "465"

Snowflake:
  WorkerID: 1
  DatacenterID: 1
  Epoch: 1735660800000
```

### 2) Gateway 配置（示例）

运行 `app/api/gateway/gateway.go` 时，默认读取 `etc/gateway.yaml`。

```yaml
Name: gateway.api
Host: 0.0.0.0
Port: 8888

Rpc:
  Auth:
    Etcd:
      Hosts:
        - 127.0.0.1:2379
      Key: sso.rpc
```

## 启动方式

在仓库根目录执行：

1. 启动 SSO RPC

```bash
go run app/rpc/sso/sso.go -f etc/sso.yaml
```

2. 启动 Gateway API

```bash
go run app/api/gateway/gateway.go -f etc/gateway.yaml
```

## 代码生成

定义文件修改后可使用 `Makefile`：

- `make api`：根据 `cmd/api/*.api` 生成 API 代码
- `make rpc`：根据 `cmd/rpc/*.proto` 生成 RPC 代码
- `make all`：同时生成 API/RPC
- `make swagger`：生成 swagger
- `make help`：查看全部命令

## 认证流程简述（当前）

1. `SendEmailCode`：生成 6 位验证码，写入 Redis（5 分钟），发送邮件。
2. `Register`：邮箱注册，密码加盐哈希，写入 Mongo `users`，返回双 Token。
3. `Login`：邮箱登录，状态检查 + 密码校验，返回双 Token。
4. Token：支持刷新、校验、登出接口。

## 开发约定

- API 层保持薄层：参数映射 + RPC 调用 + `CommonResp` 封装。
- 业务规则在 RPC 层实现，错误统一使用 `pkg/err`。
- 禁止提交任何真实账号/密钥配置，使用本地 `etc/*.yaml`。
- 协议改动顺序建议：`proto/api 定义` -> `生成代码` -> `logic` -> `README`。
