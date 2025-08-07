# emo_trash 项目架构文档

## 1. 文档修订历史

| 版本   | 修订人   | 修订日期       | 修订内容 |
|------|-------|------------|------|
| v1.0 | krace | 2025-08-07 | 初稿   |

## 2. 项目概述

### 简介
**emo_trash** 是一个基于 Go 语言开发的匿名情感宣泄平台，采用微服务架构设计，主打"无压力发泄"概念。平台核心功能包括：
🔥 匿名情绪宣泄
- 完全匿名发布机制 
- 支持文字/语音/图片等多种发泄形式 
- 智能脏话过滤与情绪分级系统

🤖 AI情感博弈
- 可定制AI发泄对象（可设定身份/性格/弱点）
- 多模态交互（文字对骂/语音互怼）
- 情感反射引擎（AI会记住你的发泄风格）

🌐 情绪社交网络
- 基于LBS的情绪热点地图
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
- 语言: Go 1.21+
- 框架: go-zero、go-eino
- 通信: gRPC/HTTP
- 存储: MySQL, MongoDB, Redis, Elasticsearch, Kafka
- 监控: Prometheus + Grafana

## 3. 系统架构

### 3.1 项目架构

```
.
├── app                 ## 项目代码
│   ├── api             ## api 网关
│   │   └── gateway
│   ├── rpc             ## rpc 服务
│   │   └── sso
├── cmd                 ## proto/api 定义
│   ├── api             ## api 定义
│   │   └── gateway.api
│   └── rpc             ## rpc 定义
│       └── sso.proto
├── deploy              ## 部署脚本
├── doc                 ## 文档
├── model               ## 数据库模型
├── pkg                 ## 工具包
│   ├── auth            ## 认证
│   ├── constant        ## 常量
│   ├── cors            ## 跨域
│   ├── db              ## 数据库
│   │   ├── mq          ## 消息队列
│   │   ├── no_sql      ## 非关系型数据库
│   │   ├── rdb         ## 关系型数据库
│   │   └── vdb         ## 向量数据库
│   ├── eino            ## AI模型
│   ├── email           ## 邮件
│   ├── encrypt         ## 加密
│   ├── err             ## 错误
│   ├── filter          ## 过滤
│   ├── media           ## 媒体
│   ├── nacos           ## 配置中心
│   ├── page            ## 分页
│   ├── result          ## 结果
│   ├── sensitive       ## 敏感词
│   ├── utils           ## 工具
│   └── yaml            ## yaml 配置
├── static              ## 静态资源
├── swagger             ## swagger 文档
├── test                ## 测试
├── third_party         ## 第三方protobuf依赖
└── Makefile            ## 快捷脚本
```

### 3.2 服务架构
 TODO

### 3.3 数据库架构
TODO

### 3.4 消息队列架构
TODO

### 3.5 缓存架构
TODO

### 3.6 监控架构
TODO


## 4. 系统组件

### 4.1 Makefile
```aiignore
# 项目配置  
PROJECT_NAME := emo_trash                       # 项目名称
PROJECT_PATH := github.com/krace-tx/emo_trash   # 项目路径
```
- make init     # 初始化项目环境
- make all      # 生成所有代码
- make api      # 生成 api 代码
- make rpc      # 生成 rpc 代码
- make mod      # 下载依赖
- make clean    # 清理生成的代码
- make swagger  # 生成 swagger 文档
- make help     # 查看帮助

* 注意： windows 系统需要通过 wsl2 来运行 make 命令


