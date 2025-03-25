package config

/**
 * @Author: 南宫乘风
 * @Description:
 * @File:  config.go
 * @Email: 1794748404@qq.com
 * @Date: 2025-03-17 17:18
 */

const (
	ListenAddress = "0.0.0.0:8081"
	//为了验证多集群
	Kubeconfigs    = `{"TST-1":"E:\\GitHUB_Code_Check\\VUE\\kubea-go\\config\\k8s.yaml","TST-2":"E:\\GitHUB_Code_Check\\VUE\\kubea-go\\config\\k8s.yaml"}`
	PodLogTailLine = 500
)
