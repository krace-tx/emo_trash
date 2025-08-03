package cors

import (
	"github.com/rs/cors"
	"github.com/zeromicro/go-zero/rest"
	"net/http"
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

func NewRestCors() rest.RunOption {
	return rest.WithCustomCors(nil, func(w http.ResponseWriter) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}, "*")
}
