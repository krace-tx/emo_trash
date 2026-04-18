# emo_trash 项目架构文档

## 1. 文档修订历史

| 版本   | 修订人   | 修订日期   | 修订内容 |
|--------|----------|------------|----------|
| v1.0   | krace    | 2025-08-07 | 初稿     |
| v1.1   | —        | 2026-04-18 | 同步目录结构、部署与离线镜像说明 |

## 2. 项目概述

### 简介

**emo_trash** 是一个基于 Go 语言开发的匿名情感宣泄平台，采用微服务架构设计，主打「无压力发泄」概念。平台核心功能包括：

🔥 匿名情绪宣泄

- 完全匿名发布机制
- 支持文字/语音/图片等多种发泄形式
- 智能脏话过滤与情绪分级系统

🤖 AI 情感博弈

- 可定制 AI 发泄对象（可设定身份/性格/弱点）
- 多模态交互（文字对骂/语音互怼）
- 情感反射引擎（AI 会记住你的发泄风格）

🌐 情绪社交网络

- 基于 LBS 的情绪热点地图
- 相似境遇推荐算法
- 匿名互动社区（点赞/共鸣/组团发泄）

💡 特色技术

- 情感分析实时处理
- 发泄内容生命周期管理
- 多重匿名保障机制
- 用户认证与授权（JWT/OAuth2）
- 社交互动管理
- AI 增强分析

**技术栈**：

- 语言: Go 1.24+
- 框架: go-zero、go-eino
- 通信: gRPC/HTTP
- 存储: MySQL、MongoDB、Redis、Elasticsearch 等（见 `pkg/datastore`）
- 消息: Kafka 等（见 `pkg/datastore/queue`）
- 监控: Prometheus + Grafana（规划）

## 3. 系统架构

### 3.1 项目架构

```
.
├── app                 ## 业务代码
│   ├── api             ## API 网关
│   │   └── gateway
│   └── rpc             ## RPC 服务
│       └── sso
├── cmd                 ## proto / api 定义
│   ├── api
│   │   └── gateway.api
│   └── rpc
│       └── sso.proto
├── deploy              ## 部署与容器
│   ├── emo_trash       ## 基础中间件栈（Redis / MySQL / MongoDB / etcd）
│   │   ├── docker-compose.yaml          ## 在线：可从仓库拉镜像
│   │   └── docker-compose.offline.yaml ## 离线：仅使用本机已 load 的镜像（pull_policy: never）
│   ├── make            ## Go 开发容器（Dockerfile + compose）
│   │   ├── docker-compose.yaml
│   │   └── docker-compose.offline.yaml ## 离线：无 build，仅 make-go-dev 镜像
│   └── offline         ## 镜像离线打包 / 导入
│       ├── export-offline.bat / .sh    ## 按镜像分别 docker save 到 images/*.tar
│       ├── import-offline.bat / .sh    ## 依次 docker load（与 save 成对使用）
│       └── images/     ## 导出的 tar 存放目录（见 .gitignore）
├── doc                 ## 文档
├── model               ## 数据库模型
├── pkg                 ## 公共库
│   ├── auth
│   ├── constant
│   ├── cors
│   ├── datastore       ## 存储与消息适配（原 pkg/db 已拆分为语义化子包）
│   │   ├── sqlstore    ## 关系型（GORM / MySQL 等）
│   │   ├── redis       ## Redis 客户端封装
│   │   ├── mongo       ## MongoDB 连接与配置
│   │   ├── search      ## Elasticsearch / 向量检索等
│   │   └── queue       ## Kafka / RabbitMQ / RocketMQ 等
│   ├── eino
│   ├── email
│   ├── encrypt
│   ├── err
│   ├── filter
│   ├── interceptor
│   ├── media
│   ├── nacos
│   ├── oss
│   ├── page
│   ├── result
│   ├── sensitive
│   ├── snowflake
│   ├── utils
│   └── yaml
├── static
├── swagger
├── test
├── third_party
├── Makefile
└── README.md
```

### 3.2 离线镜像与 Compose

1. **导出**：在仓库根目录或 `deploy/offline` 下执行 `export-offline.bat`（Windows）或 `export-offline.sh`（Git Bash / Linux），在 `deploy/offline/images/` 下为每个镜像生成独立 `.tar`（`docker save`）。
2. **导入**：将 `images/` 拷贝到离线机后执行 `import-offline.bat` / `import-offline.sh`，使用 `docker load -i` 依次加载（`docker save` 的归档需用 `load`，不要用 `import`）。
3. **启动**：中间件栈使用 `deploy/emo_trash/docker-compose.offline.yaml`；开发容器使用 `deploy/make/docker-compose.offline.yaml`（`pull_policy: never`，不拉取、不构建）。

### 3.3 服务架构

TODO

### 3.4 数据库架构

TODO

### 3.5 消息队列架构

TODO

### 3.6 缓存架构

TODO

### 3.7 监控架构

TODO

## 4. 系统组件

### 4.1 Makefile

```makefile
# 项目配置
PROJECT_NAME := emo_trash                       # 项目名称
PROJECT_PATH := github.com/krace-tx/emo_trash   # 项目路径
```

- `make init` — 初始化项目环境
- `make all` — 生成所有代码
- `make api` — 生成 API 代码
- `make rpc` — 生成 RPC 代码
- `make mod` — 下载依赖
- `make clean` — 清理生成的代码
- `make swagger` — 生成 swagger 文档
- `make help` — 查看帮助

*注意：Windows 若未安装 Make，可通过 WSL2 或 MSYS2 运行 `make`。*
