package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"kubea-go/config"

	"github.com/aryming/logger"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

/**
 * @Author: 南宫乘风
 * @Description:
 * @File:  pod.go
 * @Email: 1794748404@qq.com
 * @Date: 2025-03-20 15:42
 */

var Pod pod

type pod struct {
}

// PodsResp 定义列表的返回内容，Items是pod元素列表，Total为pod元素数量
type PodsResp struct {
	Items []corev1.Pod `json:"items"`
	Total int          `json:"total"`
}

// GetPods 获取pod列表，支持过滤和分页,排序
func (p *pod) GetPods(client *kubernetes.Clientset, filterName, namespace string, limit, page int) (*PodsResp, error) {
	// 获取podList类型的pod列表
	podList, err := client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(errors.New("获取Pod列表失败, " + err.Error()))
		return nil, errors.New("获取Pod列表失败, " + err.Error())
	}
	// 实例化dataSelector对象
	selectableData := &dataSelector{
		GenericDateSelect: p.toCells(podList.Items),
		dataSelectQuery: &DataSelectQuery{
			FilterQuery: &FilterQuery{Name: filterName},
			PaginationQuery: &PaginationQuery{
				Limit: limit,
				Page:  page,
			},
		},
	}
	//先过滤
	filtered := selectableData.Filter()
	total := len(filtered.GenericDateSelect)
	data := filtered.Sort().Paginate()
	//将[]DataCell类型的pod列表转为v1.pod列表
	pods := p.fromCells(data.GenericDateSelect)
	return &PodsResp{Items: pods, Total: total}, nil
}

// GetPodDetail 获取pod详情
func (p *pod) GetPodDetail(client *kubernetes.Clientset, namespace, podName string) (*corev1.Pod, error) {
	pod, err := client.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		logger.Error(errors.New("获取Pod详情失败, " + err.Error()))
		return nil, errors.New("获取Pod详情失败, " + err.Error())
	}
	return pod, nil
}

// DeletePod 删除POD
func (p *pod) DeletePod(client *kubernetes.Clientset, namespace, podName string) error {
	// 删除pod
	err := client.CoreV1().Pods(namespace).Delete(context.TODO(), podName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error(errors.New("删除Pod失败, " + err.Error()))
		return errors.New("删除Pod失败, " + err.Error())
	}
	return nil
}

// UpdatePod 更新POD   content参数是请求中传入的pod对象的json数据
func (p *pod) UpdatePod(client *kubernetes.Clientset, namespace, podName, content string) error {
	var pod = &corev1.Pod{}
	// 将content参数的json数据解析到pod对象中
	err := json.Unmarshal([]byte(content), pod)
	if err != nil {
		logger.Error(errors.New("反序列化失败, " + err.Error()))
		return errors.New("反序列化失败, " + err.Error())
	}
	// 更新pod
	_, err = client.CoreV1().Pods(namespace).Update(context.TODO(), pod, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(errors.New("更新Pod失败, " + err.Error()))
		return errors.New("更新Pod失败, " + err.Error())
	}
	return nil
}

// GetPodContainer 获取pod容器
func (p *pod) GetPodContainer(client *kubernetes.Clientset, namespace, podName string) (containers []string, err error) {
	// 获取pod详情
	pod, err := p.GetPodDetail(client, namespace, podName)
	if err != nil {
		logger.Error(errors.New("获取Pod详情失败, " + err.Error()))
		return nil, errors.New("获取Pod详情失败, " + err.Error())
	}
	// 从pod详情中获取容器列表
	for _, container := range pod.Spec.Containers {
		containers = append(containers, container.Name)
	}
	return containers, nil
}

// GetPodLog 获取pod内容器日志
func (p *pod) GetPodLog(client *kubernetes.Clientset, namespace, podName, containerName string) (log string, err error) {
	//设置日志的配置，容器名、tail的行数
	lineLimit := int64(config.PodLogTailLine)
	options := &corev1.PodLogOptions{
		Container: containerName,
		TailLines: &lineLimit,
	}
	// 获取request实例
	req := client.CoreV1().Pods(namespace).GetLogs(podName, options)
	// 发起request请求，返回一个io.ReadCloser类型（等同于response.body）
	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		logger.Error(errors.New("获取Pod日志失败, " + err.Error()))
		return "", errors.New("获取Pod日志失败, " + err.Error())
	}
	defer podLogs.Close()
	//将response body写入到缓冲区，目的是为了转成string返回
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		logger.Error(errors.New("复制PodLog失败, " + err.Error()))
		return "", errors.New("复制PodLog失败, " + err.Error())
	}
	return buf.String(), nil
}

// toCells 方法用于将pod类型数组，转换成DataCell类型数组
func (p *pod) toCells(std []corev1.Pod) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = podCell(std[i])
	}
	return cells
}

// fromCells 方法用于将DataCell类型数组，转换成pod类型数组
func (p *pod) fromCells(cells []DataCell) []corev1.Pod {
	pods := make([]corev1.Pod, len(cells))
	for i := range cells {
		//cells[i].(podCell)就使用到了断言,断言后转换成了podCell类型，然后又转换成了Pod类型
		pods[i] = corev1.Pod(cells[i].(podCell))
	}
	return pods
}
