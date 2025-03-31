package service

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/aryming/logger"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

/**
 * @Author: 南宫乘风
 * @Description:
 * @File:  daemonset.go
 * @Email: 1794748404@qq.com
 * @Date: 2025-03-31 18:13
 */

var Daemonset daemonset

type daemonset struct{}

type DaemonsetsResp struct {
	Items []appsv1.DaemonSet `json:"items"`
	Total int                `json:"total"`
}

func (d *daemonset) GetDaemonSets(client *kubernetes.Clientset, filterName string, namespace string, limit int, page int) (*DaemonsetsResp, error) {
	//获取daemonSetList类型的daemonSet列表
	daemonSetList, err := client.AppsV1().DaemonSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(errors.New("获取DaemonSet列表失败, " + err.Error()))
		return nil, errors.New("获取DaemonSet列表失败, " + err.Error())
	}
	//将daemonSetList中的daemonSet列表(Items)，放进dataselector对象中，进行排序
	selectableData := &dataSelector{
		GenericDateSelect: d.toCells(daemonSetList.Items),
		dataSelectQuery: &DataSelectQuery{
			FilterQuery: &FilterQuery{Name: filterName},
			PaginationQuery: &PaginationQuery{
				Limit: limit,
				Page:  page,
			},
		},
	}

	filtered := selectableData.Filter()
	total := len(filtered.GenericDateSelect)
	data := filtered.Sort().Paginate()
	//将[]DataCell类型的daemonset列表转为v1.daemonset列表
	daemonSets := d.fromCells(data.GenericDateSelect)
	return &DaemonsetsResp{
		Items: daemonSets,
		Total: total,
	}, nil
}

// GetDaemonSetDetail 获取daemonset详情
func (d *daemonset) GetDaemonSetDetail(client *kubernetes.Clientset, namespace, name string) (*appsv1.DaemonSet, error) {
	daemonSet, err := client.AppsV1().DaemonSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logger.Error(errors.New("获取DaemonSet详情失败, " + err.Error()))
		return nil, errors.New("获取DaemonSet详情失败, " + err.Error())
	}
	return daemonSet, nil
}

// DeleteDaemonSet 删除daemonset
func (d *daemonset) DeleteDaemonSet(client *kubernetes.Clientset, namespace, name string) error {
	err := client.AppsV1().DaemonSets(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		logger.Error(errors.New("删除DaemonSet失败, " + err.Error()))
		return errors.New("删除DaemonSet失败, " + err.Error())
	}
	return nil
}

// UpdateDaemonSet 更新daemonset
func (d *daemonset) UpdateDaemonSet(client *kubernetes.Clientset, namespace, content string) error {
	var daemonSet = &appsv1.DaemonSet{}
	err := json.Unmarshal([]byte(content), daemonSet)
	if err != nil {
		logger.Error(errors.New("反序列化失败, " + err.Error()))
		return errors.New("反序列化失败, " + err.Error())
	}
	_, err = client.AppsV1().DaemonSets(namespace).Update(context.TODO(), daemonSet, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(errors.New("更新DaemonSet失败, " + err.Error()))
		return errors.New("更新DaemonSet失败, " + err.Error())
	}
	return nil
}

// toCells 将 DaemonSet 列表转换为 DataCell 列表。
// 参数 std: 一组 DaemonSet 对象。
// 返回值: 一个 DataCell 类型的切片，每个元素包含一个转换后的 DaemonSet。
func (d *daemonset) toCells(std []appsv1.DaemonSet) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = daemonSetCell(std[i])
	}
	return cells
}

// fromCells 将 DataCell 列表转换回 DaemonSet 列表。
// 参数 cells: 一组 DataCell 对象，每个都包含一个 DaemonSet。
// 返回值: 一个 appsv1.DaemonSet 类型的切片，包含转换后的 DaemonSet 对象。
func (d *daemonset) fromCells(cells []DataCell) []appsv1.DaemonSet {
	daemonSets := make([]appsv1.DaemonSet, len(cells))
	for i := range cells {
		daemonSets[i] = appsv1.DaemonSet(cells[i].(daemonSetCell))
	}
	return daemonSets
}
