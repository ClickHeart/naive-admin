package server

import (
	"context"
	"fmt"
	"naive-admin/internal/router"
	"naive-admin/pkg/config"
	"naive-admin/pkg/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func Init() (app *http.Server) {
	r := gin.Default()
	router.Init(r)
	app = &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Conf.Server.Http.Port),
		Handler: r,
	}
	go func() {
		if err := app.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	return
}

func Exit(app *http.Server) {
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 会发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 会发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutdown Server ...")
	// 设置 5 秒的超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Info("Server exiting")
}
