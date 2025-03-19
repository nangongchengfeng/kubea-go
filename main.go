package main

import (
	"context"
	"kubea-go/config"
	"kubea-go/controller"
	"kubea-go/service"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/aryming/logger"

	"github.com/gin-gonic/gin"
)

/**
 * @Author: 南宫乘风
 * @Description:
 * @File:  main.go
 * @Email: 1794748404@qq.com
 * @Date: 2025-03-17 14:49
 */
func main() {
	logger.SetLogger("config/log.json")
	// 初始化路由
	r := gin.Default()
	// 初始化K8S客户端
	service.K8s.Init()
	controller.Router.InitApiRouter(r)
	// 启动服务
	srv := &http.Server{
		Addr:    config.ListenAddress,
		Handler: r,
	}
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("listen: %s\n", err)
		}
	}()
	//优雅关闭
	// 声明一个系统信号的channel，并监听系统信号，如果没有信号，就一直阻塞，如果有，就继续执行。当接收到中断信号时，执行cancel()
	// 创建一个用于接收OS信号的通道
	quit := make(chan os.Signal)
	// 配置信号通知，将OS中断信号通知到quit通道
	signal.Notify(quit, os.Interrupt)
	// 阻塞等待，直到从quit通道接收到信号
	<-quit

	// 设置上下文对象ctx，带有5秒的超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// 当函数返回时，调用cancel以取消上下文并释放资源
	defer cancel()

	// 尝试优雅关闭GIN服务器
	if err := srv.Shutdown(ctx); err != nil {
		// 如果关闭失败，记录致命错误
		logger.Error("Gin Server Shutdown:", err)
	}
	// 记录服务器退出信息
	logger.Info("Gin Server exiting")
}
