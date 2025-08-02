package filter

import (
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"time"
)

// RouteLogFilter 记录路由访问日志
func RouteLogFilter(r *http.Request) (*http.Request, error) {
	start := time.Now() // 记录开始时间
	// 在这里可以执行任何路由日志逻辑，例如保存到数据库或其他地方
	logx.Infof("Accessed: %s %s, duration: %s", r.Method, r.URL.Path, time.Since(start))
	return r, nil
}
