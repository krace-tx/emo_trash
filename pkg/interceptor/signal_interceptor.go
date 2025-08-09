package interceptor

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/krace-tx/emo_trash/pkg/db/rdb"
)

func SignalInterceptor() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		// 阻塞等待信号
		sig := <-sigChan
		log.Printf("收到关闭信号: %v，开始关闭...", sig)

		// 1. 关闭数据库连接
		if err := rdb.Close(); err != nil {
			log.Printf("数据库连接关闭失败: %v", err)
		} else {
			log.Println("数据库连接已关闭")
		}

		log.Println("所有资源已清理，服务退出")
		os.Exit(0)
	}()
}
