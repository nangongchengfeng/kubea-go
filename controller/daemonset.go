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
 * @File:  daemonset.go
 * @Email: 1794748404@qq.com
 * @Date: 2025-03-31 18:29
 */

var DaemonSet daemonSet

type daemonSet struct{}

// GetDaemonSets 获取daemonset列表，支持过滤、排序、分页
func (d *daemonSet) GetDaemonSets(ctx *gin.Context) {
	params := new(struct {
		FilterName string `form:"filter_name"`
		Namespace  string `form:"namespace"`
		Page       int    `form:"page"`
		Limit      int    `form:"limit"`
		Cluster    string `form:"cluster"`
	})
	if err := ctx.Bind(params); err != nil {
		logger.Error("Bind请求参数失败," + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		logger.Error("获取k8s连接失败," + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	daemonsets, err := service.Daemonset.GetDaemonSets(client, params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		logger.Error("获取daemonset列表失败," + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取daemonset列表成功",
		"data": daemonsets,
	})
}

// GetDaemonSetDetail 获取daemonset详情
func (d *daemonSet) GetDaemonSetDetail(ctx *gin.Context) {
	params := new(struct {
		Namespace     string `form:"namespace"`
		DaemonSetName string `form:"daemonset_name"`
		Cluster       string `form:"cluster"`
	})
	if err := ctx.Bind(params); err != nil {
		logger.Error("Bind请求参数失败," + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		logger.Error("获取k8s连接失败," + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	daemonset, err := service.Daemonset.GetDaemonSetDetail(client, params.Namespace, params.DaemonSetName)
	if err != nil {
		logger.Error("获取daemonset详情失败," + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取daemonset详情成功",
		"data": daemonset,
	})
}

// DeleteDaemonSet 删除daemonset
func (d *daemonSet) DeleteDaemonSet(ctx *gin.Context) {
	params := new(struct {
		Namespace     string `form:"namespace"`
		DaemonSetName string `form:"daemonset_name"`
		Cluster       string `form:"cluster"`
	})
	//DELETE请求，绑定参数方法改为ctx.ShouldBindJSON
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error("Bind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		logger.Error("获取k8s连接失败," + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	err = service.Daemonset.DeleteDaemonSet(client, params.Namespace, params.DaemonSetName)
	if err != nil {
		logger.Error("删除daemonset失败," + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return

	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "删除daemonset成功",
		"data": nil,
	})
}

// 更新daemonset
func (d *daemonSet) UpdateDaemonSet(ctx *gin.Context) {
	params := new(struct {
		Namespace string `json:"namespace"`
		Content   string `json:"content"`
		Cluster   string `json:"cluster"`
	})
	//PUT请求，绑定参数方法改为ctx.ShouldBindJSON
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error("Bind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		logger.Error("获取k8s连接失败," + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	err = service.Daemonset.UpdateDaemonSet(client, params.Namespace, params.Content)
	if err != nil {
		logger.Error("更新daemonset失败," + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "更新daemonset成功",
		"data": nil,
	})
}
