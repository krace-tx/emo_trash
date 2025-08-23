package errx

const (
	// ===== 模块编号 =====
	ModuleAuth   = 10 // 通用认证
	ModuleSSO    = 11 // 单点登录
	ModuleOAuth  = 12 // 第三方 OAuth
	ModuleCAS    = 13 // CAS 协议
	ModuleDB     = 20 // 数据库
	ModuleUser   = 30 // 用户业务
	ModuleSystem = 50 // 系统
	ModuleCache  = 60 // 缓存
)

// NewErr 统一创建错误码
func NewErr(module, code int, msg string) *Err {
	return New(uint32(module*1000+code), msg)
}

var (
	// ====== 通用认证模块 (10xxx) ======
	ErrAuthUnauthorized          = NewErr(ModuleAuth, 1, "用户未授权")        // 10001
	ErrAuthForbidden             = NewErr(ModuleAuth, 2, "令牌过期或无效")      // 10002
	ErrAuthTokenRequired         = NewErr(ModuleAuth, 3, "令牌不能为空")       // 10003
	ErrAuthLoginExpired          = NewErr(ModuleAuth, 4, "登录状态已过期")      // 10004
	ErrAuthRateLimit             = NewErr(ModuleAuth, 5, "接口限流，请稍后再试")   // 10005
	ErrAuthPasswordIncorrect     = NewErr(ModuleAuth, 6, "账号或密码错误")      // 10006
	ErrAuthPasswordVerifyError   = NewErr(ModuleAuth, 7, "密码验证异常")       // 10007
	ErrAuthGenAccessTokenFail    = NewErr(ModuleAuth, 8, "生成访问令牌失败")     // 10008
	ErrAuthGenRefreshTokenFail   = NewErr(ModuleAuth, 9, "生成刷新令牌失败")     // 10009
	ErrAuthSSOCheckFail          = NewErr(ModuleAuth, 10, "单点登录状态检查失败")  // 10010
	ErrAuthTokenBlacklistFail    = NewErr(ModuleAuth, 11, "旧令牌加入黑名单失败")  // 10011
	ErrAuthSaveLoginRecordFail   = NewErr(ModuleAuth, 12, "保存登录状态失败")    // 10012
	ErrAuthMobileFormat          = NewErr(ModuleAuth, 13, "手机号格式不正确")    // 10013
	ErrAuthPasswordFormat        = NewErr(ModuleAuth, 14, "密码格式不符合要求")   // 10014
	ErrAuthSmsCodeExpired        = NewErr(ModuleAuth, 15, "短信验证码已过期")    // 10015
	ErrAuthSmsCodeInvalid        = NewErr(ModuleAuth, 16, "短信验证码无效或已过期") // 10016
	ErrAuthMobileExists          = NewErr(ModuleAuth, 17, "手机号已注册")      // 10017
	ErrAuthGenIDFailed           = NewErr(ModuleAuth, 18, "生成用户ID失败")    // 10017
	ErrAuthPwdEncryptFail        = NewErr(ModuleAuth, 19, "密码加密失败")      // 10019
	ErrAuthCreateUserAuthFail    = NewErr(ModuleAuth, 20, "创建用户认证记录失败")  // 10020
	ErrAuthCreateUserProfileFail = NewErr(ModuleAuth, 21, "创建用户资料失败")    // 10021
	ErrAuthRegisterFailed        = NewErr(ModuleAuth, 22, "用户注册失败")      // 10022

	// ====== SSO 单点登录模块 (11xxx) ======
	ErrSSOTicketInvalid   = NewErr(ModuleSSO, 1, "票据无效")      // 11001
	ErrSSOTicketExpired   = NewErr(ModuleSSO, 2, "票据已过期")     // 11002
	ErrSSOServiceInvalid  = NewErr(ModuleSSO, 3, "非法服务标识")    // 11003
	ErrSSOTokenExchange   = NewErr(ModuleSSO, 4, "令牌交换失败")    // 11004
	ErrSSOSessionNotFound = NewErr(ModuleSSO, 5, "会话不存在或已失效") // 11005
	ErrSSOLogoutFailed    = NewErr(ModuleSSO, 6, "单点登出失败")    // 11006

	// ====== OAuth2 模块 (12xxx) ======
	ErrOAuthCodeInvalid   = NewErr(ModuleOAuth, 1, "授权码无效")   // 12001
	ErrOAuthCodeExpired   = NewErr(ModuleOAuth, 2, "授权码已过期")  // 12002
	ErrOAuthClientInvalid = NewErr(ModuleOAuth, 3, "非法客户端")   // 12003
	ErrOAuthScopeDenied   = NewErr(ModuleOAuth, 4, "授权范围不足")  // 12004
	ErrOAuthTokenInvalid  = NewErr(ModuleOAuth, 5, "第三方令牌无效") // 12005

	// ====== CAS 模块 (13xxx) ======
	ErrCASTicketInvalid = NewErr(ModuleCAS, 1, "CAS 票据无效")    // 13001
	ErrCASServiceDenied = NewErr(ModuleCAS, 2, "服务未注册或被禁止访问") // 13002

	// ====== 数据库模块 (20xxx) ======
	ErrDBConnectFailed = NewErr(ModuleDB, 0, "数据库连接失败") // 20000
	ErrDBQueryFailed   = NewErr(ModuleDB, 1, "数据库查询错误") // 20001

	// ====== 用户模块 (30xxx) ======
	ErrUserNotFound       = NewErr(ModuleUser, 1, "用户不存在")      // 30001
	ErrUserDisabled       = NewErr(ModuleUser, 2, "用户已被禁用")     // 30002
	ErrUserProfileFailed  = NewErr(ModuleUser, 3, "获取用户资料失败")   // 30003
	ErrUserProfileInvalid = NewErr(ModuleUser, 4, "用户资料不完整")    // 30004
	ErrUserSearchEmpty    = NewErr(ModuleUser, 5, "未找到符合条件的用户") // 30005

	// ====== 系统模块 (50xxx) ======
	ErrSystemInternal   = NewErr(ModuleSystem, 0, "系统内部错误") // 50000
	ErrSystemArgInvalid = NewErr(ModuleSystem, 1, "参数错误")   // 50001

	// ====== 缓存模块 (60xxx) ======
	ErrCacheNoData = NewErr(ModuleCache, 0, "缓存中无数据") // 60000
)
