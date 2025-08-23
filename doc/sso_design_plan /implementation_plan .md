## SSO 服务实现文档

## 1. 总览

本方案实现统一认证与账号绑定能力，覆盖：

* 登录：手机号+短信、账号+密码、第三方（微信/QQ/支付宝）
* PC 扫码登录（生成/轮询/确认）
* 会话校验、登出、刷新 Token
* 账号注册、密码重置（短信/邮箱）
* 绑定/解绑：手机号、邮箱、三方账号

### 1.1 统一名词

* **AT**：Access Token（JWT，短效）
* **RT**：Refresh Token（Opaque 随机串，长效）
* **Device Session**：同一用户在一个设备上的登录态记录

---

## 2. 架构设计

### 2.1 逻辑架构图

```mermaid
flowchart TD
  subgraph Client[Client]
    A1[App/H5]:::c
    A2[PC Web]:::c
  end
  subgraph Edge[接入层]
    GW[API Gateway(REST/gRPC-HTTP)]
  end
  subgraph SSO[SSO域]
    S1[auth RPC(go-zero)]
    S2[Session Manager]
    S3[Captcha/RateLimiter]
    S4[OAuth Adapter(WeChat/QQ/Alipay)]
  end
  subgraph Data[数据与中间件]
    R[(Redis Cluster)]
    DB[(MySQL)]
    MQ[[Kafka 可选]]
    NAC[Nacos 10.0.6.21:8848]
    LOG1[(Mongo 分布式日志)]
    ES[(Elasticsearch 检索)]
  end
  subgraph 3rd[外部服务]
    SMS[(短信服务商)]
    MAIL[(SMTP/邮件服务)]
    OIDC[WeChat/QQ/Alipay]
  end

  A1 -- gRPC-HTTP --> GW --> S1
  A2 -- gRPC-HTTP --> GW
  S1 <--> S2
  S1 --> R
  S1 --> DB
  S1 --> S3
  S1 --> S4 --> OIDC
  S1 --> SMS
  S1 --> MAIL
  S1 -. 业务埋点 .-> MQ
  MQ --> LOG1 --> ES
  S1 <--> NAC

  classDef c fill:#eef,stroke:#88f
```

### 2.2 部署与分层

* **服务注册/配置**：Nacos；`auth` 服务以 gRPC 形式注册，暴露 gRPC 与 gRPC-HTTP。
* **数据存储**：

    * MySQL：用户主数据、凭证、绑定、会话、二维码、登录日志快照
    * Redis：验证码、限流、二维码状态、黑名单/撤销列表、分布式锁
    * Kafka→Mongo→ES（可选）：全量审计日志与检索
* **安全基线**：

    * AT=JWT(签名+Kid)，RT=随机串(数据库/Redis持久化，旋转/一次性)
    * 密码 PBKDF2/bcrypt/scrypt 存储
    * 短信/邮箱验证码风控（频率/IP/设备/黑名单/同号多次）
    * 重要操作（解绑/重置）强认证（二次校验）

---

## 3. 数据模型（核心表）

> 仅列关键字段，真实表可补齐通用审计列（id/biz\_id/created\_at/updated\_at/operator 等）。

* `users`

    * `id`(PK)、`mobile`(UK, null可空)、`email`(UK, 可空)、`nickname`、`status`、`register_ip`、`register_device`
* `user_credentials`

    * `user_id`(UK)
    * `password_hash`、`password_algo`、`password_updated_at`
* `oauth_accounts`

    * `user_id`、`platform`(wechat/qq/alipay)、`open_id`(UK平台内唯一)、`union_id`、`bind_at`
* `sessions`

    * `sid`(PK)、`user_id`、`device_type`、`device_id`、`agent`、`login_ip`、`at_expires`、`rt_hash`、`rt_expires`、`status`
* `qrcodes`

    * `qid`(PK)、`device_id`(PC) 、`status`(WAITING/SCANNED/CONFIRMED/EXPIRED)、`image_url`、`created_at`、`expires_at`、`confirmed_sid`
* `verify_codes`

    * `channel`(sms/email)、`scene`、`target`、`code_hash`、`expire_at`、`send_ip`、`send_count`
* `login_logs`

    * `user_id`、`scene`、`success`、`ip`、`device`、`detail`

**Redis Key 规划**（示例）：

* `sms:{scene}:{mobile}` → 验证码+TTL
* `email:{scene}:{email}` → 验证码+TTL
* `qr:{qid}` → {status,pc\_device\_id, image\_url, sid?}
* `rt:{sid}` → RT 哈希+过期
* `jwt:revoked:{jti}` → 黑名单（可选）
* 限流：`rl:sms:{mobile}`、`rl:login:{ip}`

---

## 4. Token 策略

* **AT**：JWT，`sub=user_id`，包含 `device_type/device_id/sid/jti`；有效期：15\~60min（可配）。
* **RT**：随机串，只保存在服务端（`sessions.rt_hash`），有效期：7\~30d（可配）；每次刷新**旋转**（旧RT立刻失效）。
* **登出**：置 `sessions.status=revoked`，清理 Redis 中 `rt:{sid}`，可选将 `jti` 加入黑名单。
* **单设备并发策略**（可配）：

    * 允许多设备：同一用户多 `sid` 并存；
    * 仅一台设备：新登录顶掉旧会话（写 `login_logs.message`："账号已在其他设备登录"）。

---

## 5. 接口流程图（逐 RPC）

> 说明：为便于阅读，图内省略异常分支的详细错误码，均在实现章节给出。`validate.proto` 的字段校验在进入逻辑前完成。

### 5.1 LoginByMobile（短信登录）

```mermaid
sequenceDiagram
  participant C as Client
  participant A as auth.LoginByMobile
  participant R as Redis
  participant D as MySQL
  C->>A: mobile, sms_code, device
  A->>R: GET sms:login:mobile
  alt code valid & not expired
    A->>D: UPSERT users(user by mobile)
    A->>D: INSERT sessions(sid, device, ip, rt_hash)
    A->>C: LoginResp(AT,RT,expires)
  else invalid
    A-->>C: error(验证码错误/过期)
  end
```

### 5.2 LoginByPassword（账号密码登录）

```mermaid
sequenceDiagram
  participant C as Client
  participant A as auth.LoginByPassword
  participant D as MySQL
  C->>A: account,password,device
  A->>D: SELECT users, user_credentials
  alt password ok
    A->>D: INSERT sessions & RT
    A-->>C: LoginResp
  else
    A-->>C: error(账号或密码错误)
  end
```

### 5.3 LoginByThirdParty（三方登录）

```mermaid
sequenceDiagram
  participant C as Client
  participant A as auth.LoginByThirdParty
  participant O as OAuth Adapter
  participant D as MySQL
  C->>A: platform, open_id, union_id
  A->>O: Verify(open_id/union_id)
  alt valid
    A->>D: SELECT/UPSERT oauth_accounts -> user
    A->>D: INSERT sessions & RT
    A-->>C: LoginResp
  else
    A-->>C: error(三方认证失败)
  end
```

### 5.4 GenerateQrcode（生成二维码，PC）

```mermaid
flowchart TD
  I[PC 设备ID, IP] -->|RPC| S[auth.GenerateQrcode]
  S --> R[(Redis set qr:{qid},TTL=120s)]
  S --> U[CDN生成 image_url]
  S --> O[返回 qid,image_url]
```

### 5.5 CheckQrcodeStatus（轮询二维码状态，PC）

```mermaid
flowchart TD
  PC -->|qid| Auth
  Auth --> R[(Redis get qr:{qid})]
  R -->|WAITING/SCANNED/CONFIRMED/EXPIRED| PC
  note right of PC: 当状态=CONFIRMED，附带AT/RT返回
```

### 5.6 ConfirmQrcodeLogin（手机确认）

```mermaid
sequenceDiagram
  participant M as Mobile(App已登录)
  participant A as auth.ConfirmQrcodeLogin
  participant R as Redis
  participant D as MySQL
  M->>A: qid, app_token
  A->>A: VerifyToken(app_token)
  alt token ok & qr WAITING/SCANNED
    A->>D: INSERT sessions for PC & RT
    A->>R: SET qr:{qid} = CONFIRMED + sid + AT/RT
    A-->>M: LoginResp(PC会话信息)
  else
    A-->>M: error
  end
```

### 5.7 VerifyToken（会话校验）

```mermaid
flowchart LR
  C[Client AT] --> A[auth.VerifyToken]
  A --> J[JWT验签+过期校验]
  J -->|ok| S[查询 sessions.sid 状态]
  S -->|active| R[返回 user_id/device]
  S -->|revoked| E[error(无效会话)]
```

### 5.8 Logout（登出）

```mermaid
flowchart LR
  C[Client token] --> A[auth.Logout]
  A --> D[(MySQL update sessions.status=revoked)]
  A --> R[(DEL rt:{sid})]
  A --> O[返回 success,true]
```

### 5.9 Register（注册）

```mermaid
sequenceDiagram
  participant C as Client
  participant A as auth.Register
  participant R as Redis
  participant D as MySQL
  C->>A: mobile,password,sms_code
  A->>R: 校验验证码
  alt 通过
    A->>D: INSERT users + credentials
    A->>D: INSERT sessions & RT
    A-->>C: RegisterResp(AT/RT)
  else
    A-->>C: error
  end
```

### 5.10 SendSmsCode / SendEmailCode（验证码发送）

```mermaid
flowchart TD
  C --> A[SendSmsCode/SendEmailCode]
  A --> G[风控: 限流/黑名单/阈值]
  G -->|ok| P[供应商发送]
  P --> R[(Redis set code + TTL)]
  R --> O[返回 expire_seconds]
```

### 5.11 ResetPassword / ResetPasswordByEmail（重置密码）

```mermaid
sequenceDiagram
  C->>A: identifier + code + new_password
  A->>R: 校验验证码
  alt ok
    A->>D: 更新 password_hash
    A->>D: 失效该用户所有sessions(status=revoked)
    A-->>C: success
  else
    A-->>C: error
  end
```

### 5.12 RefreshToken（刷新 Token）

```mermaid
sequenceDiagram
  C->>A: refresh_token
  A->>D: 查询 sessions by rt_hash
  alt valid & 未过期 & status=active
    A->>D: 旋转RT(更新rt_hash,rt_expires)
    A-->>C: 新AT/RT
  else
    A-->>C: error(刷新凭证无效)
  end
```

### 5.13 绑定/解绑（Mobile/Email/ThirdParty）

```mermaid
flowchart TD
  subgraph 绑定
  X1[VerifyToken 强认证] --> X2[校验验证码/三方凭证]
  X2 --> X3[唯一性检查(手机号/邮箱/三方唯一约束)]
  X3 --> X4[写入 users/oauth_accounts]
  X4 --> X5[返回 success]
  end
  subgraph 解绑
  Y1[VerifyToken 强认证] --> Y2[校验验证码/三方凭证]
  Y2 --> Y3[限制解绑: 至少保留一种登录方式]
  Y3 --> Y4[删除绑定记录]
  Y4 --> Y5[返回 success]
  end
```

---

## 6. 实现方案（go-zero）

### 6.1 目录结构（示例）

```
internal/
  config/
    config.go
  svc/
    service_context.go
  middleware/
    authmw.go            # AT 验签，可选
  logic/
    login/
      loginbymobilelogic.go
      loginbypasswordlogic.go
      loginbythirdpartylogic.go
    qrcode/
      generateqrcodelogic.go
      checkqrcodestatuslogic.go
      confirmqrcodeloginlogic.go
    session/
      verifytokenlogic.go
      logoutlogic.go
      refreshtokenlogic.go
    account/
      registerlogic.go
      resetpasswordlogic.go
      resetpasswordbyemaillogic.go
    bind/
      bindmobilelogic.go
      bindemaillogic.go
      bindthirdpartylogic.go
      unbindmobilelogic.go
      unbindemaillogic.go
      unbindthirdpartylogic.go
  handler/               # gRPC-HTTP 映射（如需）
  model/
    *.sql / *.go         # MySQL/XORM/Ent/GORM
  rpc/
    auth.proto           # 你的 proto
  pkg/
    token/jwt.go
    token/refresh.go
    crypto/password.go
    sms/provider.go
    email/provider.go
    oauth/adapter.go
    ratelimit/rl.go
    cache/keys.go
```

### 6.2 配置项（YAML 片段）

```yaml
Name: sso-auth
ListenOn: 0.0.0.0:9000
Etcd: {} # 如使用 Nacos，请配置 Nacos 注册与配置中心
Nacos:
  Addr: 10.0.6.21:8848
  Namespace: sso

JWT:
  Issuer: sso
  Kid: kid-2025-01
  PrivateKeyPem: |
    -----BEGIN PRIVATE KEY-----
    ...
    -----END PRIVATE KEY-----
  AccessTTL: 30m

RefreshToken:
  TTL: 720h   # 30d
  Rotate: true

Redis:
  Addr: redis:6379
  Type: node
  Pass: ""

MySQL:
  DSN: user:pass@tcp(mysql:3306)/sso?parseTime=true&loc=Local

SMS:
  Provider: aliyun
  Template:
    login: XXX
    register: YYY
  RateLimit:
    perMobilePerMin: 1
    perIpPerMin: 3

Email:
  SMTP: smtp.example.com:587
  User: no-reply@example.com
  Pass: "***"

QRCode:
  TTLSeconds: 120
```

### 6.3 关键流程细节

* **验证码**：

    * 生成 `n=6` 数字；Redis `SETNX` 防并发，TTL=5\~10min；场景隔离键
    * 下发前做频率与黑名单校验（IP、设备、目标）
* **密码**：`bcrypt(cost=12+)`；变更密码后，**强制登出**所有会话
* **JWT**：

    * Header：`{alg:RS256, kid}`；Payload：`sub, sid, jti, iat, exp, device_type, device_id`
    * 统一 `jti` 写入 Redis 黑名单（可选，按 `exp` 设置 TTL）
* **RT 旋转**：查询 `sessions` 进行 RT 校验，通过后生成新 RT，旧 RT 立即失效
* **二维码**：

    * `qid=base62(uuid)`；状态机：`WAITING→SCANNED→CONFIRMED/EXPIRED`
    * PC 轮询 `CheckQrcodeStatus`；App 扫码后调用 `ConfirmQrcodeLogin`
* **并发控制**：

    * 关键路径加 `device_id` 维度锁（如二维码确认/绑定解绑）
* **风控建议**：

    * 登录失败阈值触发滑块/图形验证码
    * 非常用地登录二次校验（短信/邮箱）

### 6.4 错误码（示例）

* `A0001` 参数错误（validate）
* `A0101` 账号或密码错误
* `A0102` 验证码错误或过期
* `A0103` 第三方认证失败
* `A0201` Token 无效或过期
* `A0202` 刷新凭证无效
* `A0301` 绑定冲突（已被占用）
* `A0302` 解绑限制（需保留至少一种方式）
* `A0401` 频率受限

> go-zero 可封装为 `errx.New(code, message)`，在 handler 层统一 `ResultBody`。

### 6.5 安全与合规

* 所有敏感日志脱敏（手机号、邮箱、IP 仅前/后缀）
* 密钥与证书管理（KMS/密钥轮换，`kid` 标识）
* 全链路审计：写 Kafka → Mongo（原始）→ ES（检索）；建立索引字段：`user_id/scene/success/ip/device/ts`
* 跨域/CSRF：扫码确认接口只接受 App 内已登录 AT

### 6.6 可观测性

* 指标：登录成功率、验证码到达率、RT 旋转成功率、平均认证时延、错误分布
* Tracing：`Login` → `Redis`/`MySQL`/`SMS` 外部依赖链路
* 报警：验证码失败率突增、第三方认证接口 5xx、RT 验证异常

---

## 7. 与 proto 对应的处理要点

| RPC                | 关键校验                  | 事务/幂等              | 失败回滚                 |
| ------------------ | --------------------- | ------------------ | -------------------- |
| LoginByMobile      | 验证码、风控                | 创建 `sessions`、写 RT | 删除已创建但未响应的会话         |
| LoginByPassword    | 密码 hash               | 同上                 | 同上                   |
| LoginByThirdParty  | open\_id/union\_id 唯一 | 绑定或新建用户            | 失败回滚 oauth\_accounts |
| GenerateQrcode     | 设备/IP 合法性             | qid 全局唯一           | 过期自动清理               |
| CheckQrcodeStatus  | qid 存在                | 仅查询                | -                    |
| ConfirmQrcodeLogin | App AT 有效             | PC 会话创建幂等          | 二维码状态CAS更新           |
| VerifyToken        | JWT 验签                | -                  | -                    |
| Logout             | 会话存在                  | 幂等                 | -                    |
| Register           | 验证码、密码强度              | 新建用户+会话            | 失败删除新用户              |
| SendSmsCode/Email  | 频控/黑名单                | 幂等（覆盖+刷新TTL）       | -                    |
| ResetPassword\*    | 强认证                   | 失效所有会话             | -                    |
| RefreshToken       | RT 有效&未过期             | 旋转                 | 旧 RT 失效              |
| Bind\*/Unbind\*    | 强认证                   | 唯一约束/保底策略          | 失败回滚绑定               |

---

## 8. 典型伪代码片段

```go
// RefreshTokenLogic.Handle
sess := repo.FindSessionByRTHash(hash(rt))
if sess == nil || sess.Status != Active || time.Now().After(sess.RTExpires) {
    return ErrRefreshInvalid
}
newRT := random()
repo.RotateRT(sess.Sid, hash(newRT), now+cfg.RTTTL)
newAT := jwt.Sign(JwtClaims{Sub:sess.UserID, Sid:sess.Sid, ...}, key)
return LoginResp{AccessToken:newAT, RefreshToken:newRT, ...}
```

```go
// VerifyTokenLogic.Handle
claims, err := jwt.ParseAndVerify(at, keySet)
if err != nil || time.Now().After(claims.Exp) { return ErrTokenInvalid }
if repo.IsRevoked(claims.Jti) { return ErrTokenInvalid }
sess := repo.GetSession(claims.Sid)
if sess.Status != Active { return ErrTokenInvalid }
return VerifyResp{UserId:claims.Sub, DeviceType:claims.DeviceType, DeviceId:claims.DeviceId}
```

---

## 9. 测试清单

* 单测：验证码发送/验证、密码校验、JWT 验签与过期、RT 旋转、二维码状态机
* 集成：登录全链路（含风控）、三方登录回调、扫码登录、绑定/解绑与保底策略
* 压测：登录 QPS、Redis 命中率、验证码通道吞吐、数据库热点行

---

## 10. 上线与运维

* 配置热更新：Nacos 配置监听（JWT TTL、验证码 TTL、频率阈值、单/多设备策略）
* 滚动发布：无状态扩缩容；会话/验证码依赖 Redis，节点可随时替换
* 备份：MySQL 周期备份；密钥/证书安全存储与轮换

---

## 11. 附录

* Proto 中 `validate.proto` 的规则即参数入站校验；业务校验在 Logic 层完成
* 建议引入 `OpenAPI`/`grpcui` 进行联调；对外仅开放需要的 gRPC-HTTP 路由
* 重要操作（解绑/重置）建议开启二次确认与最近登录设备通知
