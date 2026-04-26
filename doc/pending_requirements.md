# 后端待实现需求清单 (2025-04-26)

## 1. 社交与互动 (Post 域)
- [ ] **回声 (Echo) 系统**
  - `CreateComment`: 发布回声（支持治愈系语录快捷发布）。
  - `ListComments`: 分页获取帖子的回声列表。
  - `DeleteComment`: 用户删除自己的回声。
- [ ] **互动列表**
  - `ListStarredPosts`: 获取当前用户收藏的情绪列表。
  - `ListMyPosts`: 获取当前用户发布的情绪列表（区分匿名与非匿名）。

## 2. 用户中心 (SSO/User 域)
- [ ] **个人数据统计 (User Stats)** - *对应 `mine.uvue` 主卡片*
  - 提供接口返回：
    - `post_count`: 累计发布的情绪条数（"12 条情绪"）。
    - `resonance_count`: 获得的总拥抱/共鸣数（"36 人也这样"）。
    - `join_days`: 账号创建至今的天数（"7 天"）。
- [ ] **基础信息完善**
  - `UpdateUserInfo` 需确保支持 `bio` (签名/此刻状态) 和 `avatar` (头像) 的更新。

## 3. 治愈与情感服务 (Healing Service) - *对应 `mine.uvue` 共鸣卡片*
- [ ] **每日温柔 (Daily Healing/Comfort)**
  - `GetComfortMessage`: 获取今日推荐的治愈文案及动作。
  - 返回数据需包含：
    - `title`: 标题（如："你不是一个人"）。
    - `subtitle`: 副标题（如："今天也有人和你一样"）。
    - `button_text`: 按钮文案（如："给自己一点温柔"）。
    - `illustration_url`: 插画配置地址。

## 4. 基础设施 (Common)
- [ ] **媒体上传服务**
  - `UploadMedia`: 支持图片/附件上传到 Minio/OSS，并返回可访问 URL。（这里直接先将数据写到本地）
  - 需集成敏感词与图片审核机制（针对情绪树洞场景），通过AI 来审核。


