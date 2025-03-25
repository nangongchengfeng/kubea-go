package controller

import (
	"kubea-go/service"
	"net/http"

	"github.com/aryming/logger"
	"github.com/gin-gonic/gin"
)

/**
 * @Author: 南宫乘风
 * @Description:
 * @File:  pod.go.go
 * @Email: 1794748404@qq.com
 * @Date: 2025-03-25 15:37
 */

var Pod pod

type pod struct{}

//Controller中的方法入参是gin.Context，用于从上下文中获取请求参数及定义响应内容
//流程：绑定参数",调用service代码",根据调用结果响应具体内容

// GetPods 获取pod列表，支持过滤、排序、分页
func (p *pod) GetPods(c *gin.Context) {
	//匿名结构体，用于声明入参，get请求为form格式，其他请求为json格式
	params := new(
		struct {
			FilterName string `form:"filter_name"`
			Namespace  string `form:"namespace"`
			Page       int    `form:"page"`
			Limit      int    `form:"limit"`
			Cluster    string `form:"cluster"`
		})
	//绑定参数，给匿名结构体中的属性赋值，值是入参
	//	form格式使用ctx.Bind方法，json格式使用ctx.ShouldBindJSON方法
	if err := c.ShouldBind(params); err != nil {
		logger.Error("Bind请求参数失败," + err.Error())
		// ctx.JSON方法用于返回响应内容，入参是状态码和响应内容，响应内容放入gin.H的map中
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	// 获取k8s的连接方式
	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		logger.Error("获取k8s连接失败," + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	//service中的的方法通过 包名.结构体变量名.方法名 使用，serivce.Pod.GetPods()
	pods, err := service.Pod.GetPods(client, params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		logger.Error("获取pod列表失败," + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取pod列表成功",
		"data": pods,
	})
}

// GetPodDetail 获取pod详情
func (p *pod) GetPodDetail(cxt *gin.Context) {
	params := new(struct {
		Namespace string `form:"namespace"`
		PodName   string `form:"pod_name"`
		Cluster   string `form:"cluster"`
	})
	if err := cxt.ShouldBind(params); err != nil {
		logger.Error("Bind请求参数失败," + err.Error())
		cxt.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		logger.Error("获取k8s连接失败," + err.Error())
		cxt.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.Pod.GetPodDetail(client, params.Namespace, params.PodName)
	if err != nil {
		logger.Error("获取pod详情失败," + err.Error())
		cxt.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	cxt.JSON(http.StatusOK, gin.H{
		"msg":  "获取pod详情成功",
		"data": data,
	})
}

// DeletePod 删除pod
func (p *pod) DeletePod(cxt *gin.Context) {
	params := new(struct {
		Namespace string `form:"namespace"`
		PodName   string `form:"pod_name"`
		Cluster   string `form:"cluster"`
	})
	if err := cxt.ShouldBind(params); err != nil {
		logger.Error("Bind请求参数失败," + err.Error())
		cxt.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		logger.Error("获取k8s连接失败," + err.Error())
		cxt.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	if err := service.Pod.DeletePod(client, params.Namespace, params.PodName); err != nil {
		logger.Error("删除pod失败," + err.Error())
		cxt.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	cxt.JSON(http.StatusOK, gin.H{
		"msg":  "删除pod成功",
		"data": nil,
	})
}

// UpdatePod 更新pod
func (p *pod) UpdatePod(cxt *gin.Context) {
	params := new(struct {
		Namespace string `form:"namespace"`
		PodName   string `form:"pod_name"`
		Content   string `form:"content"`
		Cluster   string `form:"cluster"`
	})
	//PUT请求，绑定参数方法改为ctx.ShouldBindJSON
	if err := cxt.ShouldBindJSON(params); err != nil {
		logger.Error("Bind请求参数失败," + err.Error())
		cxt.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		logger.Error("获取k8s连接失败," + err.Error())
		cxt.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	if err := service.Pod.UpdatePod(client, params.Namespace, params.PodName, params.Content); err != nil {
		logger.Error("更新pod失败," + err.Error())
		cxt.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	cxt.JSON(http.StatusOK, gin.H{
		"msg":  "更新pod成功",
		"data": nil,
	})
}

// GetPodContainer 获取pod容器
func (p *pod) GetPodContainer(cxt *gin.Context) {
	params := new(struct {
		Namespace string `form:"namespace"`
		PodName   string `form:"pod_name"`
		Cluster   string `form:"cluster"`
	})
	// GET请求，绑定参数方法改为ctx.Bind
	if err := cxt.Bind(params); err != nil {
		logger.Error("Bind请求参数失败," + err.Error())
		cxt.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		logger.Error("获取k8s连接失败," + err.Error())
		cxt.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	containers, err := service.Pod.GetPodContainer(client, params.Namespace, params.PodName)
	if err != nil {
		logger.Error("获取pod容器失败," + err.Error())
		cxt.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	cxt.JSON(http.StatusOK, gin.H{
		"msg":  "获取pod容器成功",
		"data": containers,
	})
}

// GetPodLog 获取pod中容器日志
func (p *pod) GetPodLog(cxt *gin.Context) {
	params := new(struct {
		Namespace     string `form:"namespace"`
		PodName       string `form:"pod_name"`
		ContainerName string `form:"container_name"`
		Cluster       string `form:"cluster"`
	})
	// GET请求，绑定参数方法改为ctx.Bind
	if err := cxt.Bind(params); err != nil {
		logger.Error("Bind请求参数失败," + err.Error())
		cxt.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		logger.Error("获取k8s连接失败," + err.Error())
		cxt.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	log, err := service.Pod.GetPodLog(client, params.Namespace, params.PodName, params.ContainerName)
	if err != nil {
		logger.Error("获取pod日志失败," + err.Error())
		cxt.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	cxt.JSON(http.StatusOK, gin.H{
		"msg":  "获取pod日志成功",
		"data": log,
	})
}
