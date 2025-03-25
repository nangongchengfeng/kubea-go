package service

/**
 * @Author: 南宫乘风
 * @Description:
 * @File:  k8s_client.go
 * @Email: 1794748404@qq.com
 * @Date: 2025-03-18 17:55
 */
import (
	"encoding/json"
	"errors"
	"fmt"
	"kubea-go/config"

	"github.com/aryming/logger"
	"k8s.io/client-go/tools/clientcmd"

	"k8s.io/client-go/kubernetes"
)

var K8s k8s

type k8s struct {
	// 提供多集群Client
	ClientMap map[string]*kubernetes.Clientset
	// 提供多集群列表
	KubeConfMap map[string]string
}

// GetClient 根据集群名称获取Client
// 根据集群名称获取kubernetes客户端
func (k *k8s) GetClient(clusterName string) (*kubernetes.Clientset, error) {
	// 从ClientMap中获取指定集群名称的客户端
	client, ok := k.ClientMap[clusterName]
	// 如果不存在，则返回错误
	if !ok {
		logger.Error(fmt.Sprintf("集群%s不存在,无法获取Client\n", clusterName))
		return nil, errors.New(fmt.Sprintf("集群%s不存在,无法获取Client\n", clusterName))
	}
	// 返回客户端
	return client, nil
}

// 初始化k8s client
func (k *k8s) Init() {
	// 创建一个空的map，用于存储Kubeconfigs
	mp := make(map[string]string, 0)
	// 创建一个空的map，用于存储Kubernetes的Clientset
	k.ClientMap = make(map[string]*kubernetes.Clientset, 0)
	// 反序列化
	if err := json.Unmarshal([]byte(config.Kubeconfigs), &mp); err != nil {
		// 如果反序列化失败，则抛出异常
		panic(fmt.Sprintf("反序列化Kubeconfigs失败,%v\n", err))
	}
	// 将反序列化后的结果存储到KubeConfMap中
	k.KubeConfMap = mp
	// 初始化集群Client
	for key, value := range mp {
		// 根据Kubeconfigs中的配置，初始化集群Client
		client, err := clientcmd.BuildConfigFromFlags("", value)
		if err != nil {
			// 如果初始化失败，则抛出异常
			panic(fmt.Sprintf("初始化集群%s失败,%v\n", key, err))
		}
		clientSet, err := kubernetes.NewForConfig(client)
		if err != nil {
			// 如果初始化失败，则抛出异常
			panic(fmt.Sprintf("初始化集群%s失败,%v\n", key, err))
		}
		// 将初始化后的Clientset存储到ClientMap中
		k.ClientMap[key] = clientSet
		// 打印初始化成功的日志
		logger.Info(fmt.Sprintf("初始化集群%s成功", key))
	}
}
