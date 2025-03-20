package main

import (
	"fmt"
	"kubea-go/service"
)

/**
 * @Author: 南宫乘风
 * @Description:
 * @File:  main.go
 * @Email: 1794748404@qq.com
 * @Date: 2025-03-17 14:49
 */
func main() {

	var k8s = service.K8s
	k8s.Init()
	clientset, err := k8s.GetClient("TST-1")
	if err != nil {
		return
	}
	pods, err := service.Pod.GetPods(clientset, "", "default", 10, 1)
	if err != nil {
		return
	}
	fmt.Println(pods)
}
