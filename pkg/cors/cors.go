package middleware

import (
	"github.com/rs/cors"
)

func NewCors() *cors.Cors {
	return cors.New(cors.Options{
		// 生产环境切换允许跨域来源
		AllowedOrigins: []string{
			"http://*",
			"https://*",
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Content-Type", "Access-Control-Allow-Headers", "X-Requested-With", "Access-Control-Allow-Credentials", "User-Agent", "Content-Length", "Authorization"},
		AllowCredentials: true,  // 启用凭证
		Debug:            false, // 生产环境关闭调试模式
	})
}
