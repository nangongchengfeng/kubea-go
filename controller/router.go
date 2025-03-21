package controller

import "github.com/gin-gonic/gin"

/**
 * @Author: 南宫乘风
 * @Description:
 * @File:  router.go
 * @Email: 1794748404@qq.com
 * @Date: 2025-03-17 17:18
 */

// Router 实例化对象，可以在main.go中调用
var Router router

type router struct {
}

// InitRouter 初始化路由

func (*router) InitApiRouter(r *gin.Engine) {
	r.GET("/api", func(c *gin.Context) { c.JSON(200, gin.H{"message": "pong"}) })
}
