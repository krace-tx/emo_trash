package err

// 定义业务状态码常量
var (
	// 业务状态码 1xx
	ServerNullData               = New(101, "当前暂无数据或查询信息并不存在")
	ServerNotFoundNowTerm        = New(102, "无当前学期数据")
	ServerErrorParam             = New(103, "参数有误")
	ServerDormitoryInfoException = New(104, "学生寝室信息异常")
	ServerNoDormitoryInformation = New(105, "暂无宿舍信息")
	ServerParamNotNull           = New(106, "参数不能为空")

	// 客户端错误状态码 4xx
	AuthUnauthorized = New(401, "权限模块：用户未授权")
	AuthForbidden    = New(402, "权限模块：令牌已过期或验证不正确!")
	AuthTokenNotNull = New(403, "权限模块：令牌不能为空")
	AuthLoginExpire  = New(404, "权限模块：登录状态已过期，请重新登录")
	AuthTokenFail    = New(405, "权限模块：令牌验证失败，请尝试重新登录")
	AuthFlushFail    = New(406, "权限模块: 权限刷新失败，请重新登录")
	AuthNoPermission = New(407, "权限模块: 没有接口的访问权限!")
	AuthRequestLimit = New(408, "接口限流: 收到请求过多!请稍后再试!")
	AuthUnauth       = New(401, "权限模块：用户未授权")

	// 服务器错误状态码 5xx
	Error       = New(500, "系统内部错误")
	ErrAbnormal = New(501, "请求数据异常")
	ErrDB       = New(503, "数据库异常")
	ErrNotFound = New(504, "没有找到该数据")
	ErrVerify   = New(505, "敏感内容，审核未通过")
	ErrSearch   = New(506, "全局搜索异常")

	// Redis 相关状态码 6xx
	ErrorNotInfo = New(60000, "错误，没有解码数据")
	ErrorDecode  = New(60001, "解码失败: 参数已失效!")

	// 社区业务错误码 7xx
	CommunityNotRegister        = New(701, "用户不存在或未注册")
	CommunityUserAbnormality    = New(702, "用户信息异常")
	CommunityRegistrationFailed = New(703, "用户注册失败")
	CommunityHasRegistered      = New(704, "用户已注册")
)
