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
	// Deployment 路由服务
	deploymentGroup := r.Group(apiBasePath)
	{
		deploymentGroup.GET("/deployment", Deployment.GetDeployments)
		deploymentGroup.GET("/deployment/detail", Deployment.GetDeploymentDetail)
		deploymentGroup.PUT("/deployment/scale", Deployment.ScaleDeployment)
		deploymentGroup.DELETE("/deployment/del", Deployment.DeleteDeployment)
		deploymentGroup.PUT("/deployment/restart", Deployment.RestartDeployment)
		deploymentGroup.PUT("/deployment/update", Deployment.UpdateDeployment)
		deploymentGroup.GET("/deployment/numnp", Deployment.GetDeployNumPerNp)
		deploymentGroup.POST("/deployment/create", Deployment.CreateDeployment)
	}
	// DaemonSet 路由服务
	daemonSetGroup := r.Group(apiBasePath)
	{
		daemonSetGroup.GET("/daemonset", DaemonSet.GetDaemonSets)
		daemonSetGroup.GET("/daemonset/detail", DaemonSet.GetDaemonSetDetail)
		daemonSetGroup.DELETE("/daemonset/del", DaemonSet.DeleteDaemonSet)
		daemonSetGroup.PUT("/daemonset/update", DaemonSet.UpdateDaemonSet)
	}
	// StatefulSet 路由服务
	statefulSetGroup := r.Group(apiBasePath)
	{
		statefulSetGroup.GET("/statefulset", StatefulSet.GetStatefulSets)
		statefulSetGroup.GET("/statefulset/detail", StatefulSet.GetStatefulSetDetail)
		statefulSetGroup.DELETE("/statefulset/del", StatefulSet.DeleteStatefulSet)
		statefulSetGroup.PUT("/statefulset/update", StatefulSet.UpdateStatefulSet)
	}
}
