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

const apiBasePath = "/api/k8s"

// InitRouter 初始化路由

func (*router) InitApiRouter(r *gin.Engine) {
	r.GET("/api/ping", func(c *gin.Context) { c.JSON(200, gin.H{"message": "pong"}) })

	// Pod 路由服务
	podGroup := r.Group(apiBasePath)
	{
		podGroup.GET("/pod", Pod.GetPods)
		podGroup.GET("/pod/detail", Pod.GetPodDetail)
		podGroup.DELETE("/pod/del", Pod.DeletePod)
		podGroup.PUT("/pod/update", Pod.UpdatePod)
		podGroup.GET("/pod/container", Pod.GetPodContainer)
		podGroup.GET("/pod/log", Pod.GetPodLog)
	}
	//r.GET("/api/k8s/pod", Pod.GetPods)
	//r.GET("/api/k8s/pod/detail", Pod.GetPodDetail)
	//r.DELETE("/api/k8s/pod/del", Pod.DeletePod)
	//r.PUT("/api/k8s/pod/update", Pod.UpdatePod)
	//r.GET("/api/k8s/pod/container", Pod.GetPodContainer)
	//r.GET("/api/k8s/pod/log", Pod.GetPodLog)
}
