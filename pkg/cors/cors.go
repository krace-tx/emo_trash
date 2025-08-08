package cors

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest"
)

func NewRestCors() rest.RunOption {
	return rest.WithCustomCors(
		func(header http.Header) {
			// 允许的来源（根据需求动态设置或使用白名单）
			header.Set("Access-Control-Allow-Origin", "*") // 生产环境应替换为具体域名

			// 显式声明允许的请求头（避免使用 *）
			header.Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

			// 允许的 HTTP 方法
			header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")

			// 暴露必要的自定义头（移除冗余的 CORS 头）
			header.Set("Access-Control-Expose-Headers", "Content-Length, Content-Type")

			// 预检请求缓存时间（单位：秒）
			header.Set("Access-Control-Max-Age", "86400") // 24 小时

			// 如果需要跨域凭据（Cookie/Auth），需注释掉 Allow-Origin: * 并改用具体域名
			// header.Set("Access-Control-Allow-Credentials", "true")
		},
		nil, // 使用默认的 "不允许的请求" 处理器
		"*", // 允许的来源（可替换为白名单，如 "https://example.com"）
	)
}
