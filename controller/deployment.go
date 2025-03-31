package controller

import (
	"fmt"
	"kubea-go/service"
	"net/http"

	"github.com/aryming/logger"
	"github.com/gin-gonic/gin"
)

/**
 * @Author: 南宫乘风
 * @Description:
 * @File:  deployment.go
 * @Email: 1794748404@qq.com
 * @Date: 2025-03-27 17:16
 */

var Deployment deployment

type deployment struct{}

// GetDeployments 获取deployment列表，支持过滤、排序、分页
func (d *deployment) GetDeployments(c *gin.Context) {
	//获取参数
	params := new(struct {
		FilterName string `form:"filter_name"`
		Namespace  string `form:"namespace"`
		Page       int    `form:"page"`
		Limit      int    `form:"limit"`
		Cluster    string `form:"cluster"`
	})
	if err := c.Bind(params); err != nil {
		logger.Error("Bind请求参数失败," + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	//获取client
	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		logger.Error("获取k8s连接失败," + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.Deployment.GetDeployments(client, params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		logger.Error("获取deployment列表失败," + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取deployment列表成功",
		"data": data,
	})
}

// GetDeploymentDetail 获取deployment详情
func (d *deployment) GetDeploymentDetail(c *gin.Context) {
	params := new(struct {
		DeploymentName string `form:"deployment_name"`
		Namespace      string `form:"namespace"`
		Cluster        string `form:"cluster"`
	})
	if err := c.Bind(params); err != nil {
		logger.Error("Bind请求参数失败," + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		logger.Error("获取k8s连接失败," + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.Deployment.GetDeploymentDetail(client, params.Namespace, params.DeploymentName)
	if err != nil {
		logger.Error("获取deployment详情失败," + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取deployment详情成功",
		"data": data,
	})
}

// CreateDeployment 创建deployment
func (d *deployment) CreateDeployment(c *gin.Context) {
	var (
		deployCreate = new(service.DeployCreate)
		err          error
	)
	if err = c.ShouldBindJSON(deployCreate); err != nil {
		logger.Error("Bind请求参数失败," + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	client, err := service.K8s.GetClient(deployCreate.Cluster)
	if err != nil {
		logger.Error("获取k8s连接失败," + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	if err := service.Deployment.CreateDeployment(client, deployCreate); err != nil {
		logger.Error("创建deployment失败," + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "创建deployment成功",
		"data": nil,
	})
}

// ScaleDeployment 设置deployment副本数
func (d *deployment) ScaleDeployment(c *gin.Context) {
	params := new(struct {
		DeploymentName string `form:"deployment_name"`
		Namespace      string `form:"namespace"`
		ScaleNum       int    `form:"scale_num"`
		Cluster        string `form:"cluster"`
	})
	// PUT请求，绑定参数方法改为c.ShouldBindJSON
	if err := c.ShouldBindJSON(params); err != nil {
		logger.Error("Bind请求参数失败," + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		logger.Error("获取k8s连接失败," + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.Deployment.ScaleDeployment(client, params.DeploymentName, params.Namespace, params.ScaleNum)
	if err != nil {
		logger.Error("设置deployment副本数失败," + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "设置deployment副本数成功",
		"data": fmt.Sprintf("最新副本数: %d", data),
	})
}

// 删除deployment
func (d *deployment) DeleteDeployment(c *gin.Context) {
	params := new(struct {
		DeploymentName string `form:"deployment_name"`
		Namespace      string `form:"namespace"`
		Cluster        string `form:"cluster"`
	})
	// Delete请求，绑定参数方法改为c.ShouldBindJSON
	if err := c.ShouldBindJSON(params); err != nil {
		logger.Error("Bind请求参数失败," + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		logger.Error("获取k8s连接失败," + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	if err := service.Deployment.DeleteDeployment(client, params.DeploymentName, params.Namespace); err != nil {
		logger.Error("删除deployment失败," + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "删除deployment成功",
		"data": nil,
	})
}

// RestartDeployment 重启deployment
func (d *deployment) RestartDeployment(c *gin.Context) {
	params := new(struct {
		DeploymentName string `form:"deployment_name"`
		Namespace      string `form:"namespace"`
		Cluster        string `form:"cluster"`
	})
	// PUT 请求，绑定参数方法改为c.ShouldBindJSON
	if err := c.ShouldBindJSON(params); err != nil {
		logger.Error("Bind请求参数失败," + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		logger.Error("获取k8s连接失败," + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	if err := service.Deployment.RestartDeployment(client, params.DeploymentName, params.Namespace); err != nil {
		logger.Error("重启deployment失败," + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "重启deployment成功",
		"data": nil,
	})

}

// 更新deployment
func (d *deployment) UpdateDeployment(c *gin.Context) {
	params := new(struct {
		Namespace string `form:"namespace"`
		Content   string `form:"content"`
		Cluster   string `form:"cluster"`
	})
	// PUT 请求，绑定参数方法改为c.ShouldBindJSON
	if err := c.ShouldBindJSON(params); err != nil {
		logger.Error("Bind请求参数失败," + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		logger.Error("获取k8s连接失败," + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	if err := service.Deployment.UpdateDeployment(client, params.Namespace, params.Content); err != nil {
		logger.Error("更新deployment失败," + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "更新deployment成功",
		"data": nil,
	})
}

// 获取每个namespace的pod数量
func (d *deployment) GetDeployNumPerNp(c *gin.Context) {
	params := new(struct {
		Cluster string `form:"cluster"`
	})
	// GET 请求，绑定参数方法改为c.Bind
	if err := c.Bind(params); err != nil {
		logger.Error("Bind请求参数失败," + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		logger.Error("获取k8s连接失败," + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.Deployment.GetDeployNumPerNp(client)
	if err != nil {
		logger.Error("获取每个namespace的deployment数量失败," + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取每个namespace的deployment数量成功",
		"data": data,
	})
}
