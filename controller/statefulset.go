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
 * @File:  statefulset.go
 * @Email: 1794748404@qq.com
 * @Date: 2025-04-14 10:44
 */

var StatefulSet statefulSet

type statefulSet struct{}

func (s *statefulSet) GetStatefulSets(ctx *gin.Context) {
	param := new(struct {
		FilterName string `form:"filter_name"`
		Namespace  string `form:"namespace"`
		Page       int    `form:"page"`
		Limit      int    `form:"limit"`
		Cluster    string `form:"cluster"`
	})
	if err := ctx.Bind(param); err != nil {
		logger.Error("Bind请求参数失败," + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	client, err := service.K8s.GetClient(param.Cluster)
	if err != nil {
		logger.Error("获取k8s连接失败," + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	statefulSetsResp, err := service.Statefulset.GetStatefulSets(client, param.FilterName, param.Namespace, param.Limit, param.Page)
	if err != nil {
		logger.Error("获取statefulset失败," + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取StatefulSet列表成功",
		"data": statefulSetsResp,
	})
}

// 获取statefulset详情
func (s *statefulSet) GetStatefulSetDetail(ctx *gin.Context) {
	params := new(struct {
		StatefulSetName string `form:"statefulset_name"`
		Namespace       string `form:"namespace"`
		Cluster         string `form:"cluster"`
	})
	if err := ctx.Bind(params); err != nil {
		logger.Error("Bind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.Statefulset.GetStatefulSetDetail(client, params.Namespace, params.StatefulSetName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取StatefulSet详情成功",
		"data": data,
	})
}

// 删除statefulset
func (s *statefulSet) DeleteStatefulSet(ctx *gin.Context) {
	params := new(struct {
		StatefulSetName string `json:"statefulset_name"`
		Namespace       string `json:"namespace"`
		Cluster         string `json:"cluster"`
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	err = service.Statefulset.DeleteStatefulSet(client, params.StatefulSetName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "删除StatefulSet成功",
		"data": nil,
	})
}

// 更新statefulSet
func (s *statefulSet) UpdateStatefulSet(ctx *gin.Context) {
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	err = service.Statefulset.UpdateStatefulSet(client, params.Namespace, params.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "更新StatefulSet成功",
		"data": nil,
	})
}
