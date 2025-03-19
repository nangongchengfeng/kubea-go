package main

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"k8s.io/client-go/tools/clientcmd"
)

/**
 * @Author: 南宫乘风
 * @Description:
 * @File:  main.go
 * @Email: 1794748404@qq.com
 * @Date: 2025-03-17 14:49
 */
func main() {

	// 定义kubeconfig文件路径
	kubeconfig := "config/k8s.yaml"
	// 从kubeconfig文件中构建配置
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		// 如果构建配置失败，则抛出错误
		panic(err.Error())
	}
	// 使用配置创建kubernetes客户端
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		// 如果创建客户端失败，则抛出错误
		panic(err.Error())
	}
	// 使用客户端获取default命名空间下的所有Pod
	pods, err := client.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		// 如果获取Pod失败，则抛出错误
		panic(err.Error())
	}
	// 遍历所有Pod，并打印Pod名称
	for _, pod := range pods.Items {
		println(pod.Name)
	}
}
