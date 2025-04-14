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
 * @File:  statefulset.go
 * @Email: 1794748404@qq.com
 * @Date: 2025-04-14 10:25
 */

var Statefulset statefulset

type statefulset struct {
}

type StatefulsetsResp struct {
	Items []appsv1.StatefulSet `json:"items"`
	Total int                  `json:"total"`
}

func (s *statefulset) GetStatefulSets(client *kubernetes.Clientset, filterName string, namespace string, limit, page int) (statusfulSetsResp *StatefulsetsResp, err error) {
	//获取statefulSetList类型的statefulSet列表
	statefulSetList, err := client.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(errors.New("获取StatefulSet列表失败, " + err.Error()))
		return nil, errors.New("获取StatefulSet列表失败, " + err.Error())
	}
	//将statefulSetList中的StatefulSet列表(Items)，放进dataselector对象中，进行排序
	selectableData := &dataSelector{
		GenericDateSelect: s.toCells(statefulSetList.Items),
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
	//将[]DataCell类型的statefulset列表转为v1.statefulset列表
	statefulSets := s.fromCells(data.GenericDateSelect)
	return &StatefulsetsResp{
		Items: statefulSets,
		Total: total,
	}, nil
}

func (s *statefulset) GetStatefulSetDetail(client *kubernetes.Clientset, namespace, name string) (statefulSet *appsv1.StatefulSet, err error) {
	/*
	 * @Author: 南宫乘风
	 * @Time: 2021-06-03 15:02:01
	 * @Description: 获取指定StatefulSet的详细信息
	 * @Params:
	 *   - client *kubernetes.Clientset: Kubernetes客户端，用于连接Kubernetes集群
	 *   - namespace string: StatefulSet所在的命名空间
	 *   - name string: StatefulSet的名称
	 * @Returns:
	 *   - statefulSet *appsv1.StatefulSet: 返回获取的StatefulSet详细信息，如果获取失败则为nil
	 *   - err error: 返回错误信息，如果获取成功则为nil
	 * @Details:
	 *   - 使用Kubernetes客户端的AppsV1()方法获取StatefulSet资源
	 *   - 调用Get方法，传入命名空间、StatefulSet名称和获取选项来获取StatefulSet详情
	 *   - 如果获取失败，记录错误日志并返回自定义错误信息
	 *   - 为什么这么做：直接使用Kubernetes客户端提供的方法可以高效、便捷地获取StatefulSet的详细信息
	 */

	// 获取指定命名空间和名称的StatefulSet详情
	statefulSet, err = client.AppsV1().StatefulSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		// 如果获取失败，记录错误日志并返回自定义错误信息
		logger.Error(errors.New("获取StatefulSet详情失败, " + err.Error()))
		return nil, errors.New("获取StatefulSet详情失败, " + err.Error())
	}
	// 返回获取的StatefulSet详情
	return statefulSet, nil
}

// 删除statefulset
func (s *statefulset) DeleteStatefulSet(client *kubernetes.Clientset, statefulSetName, namespace string) (err error) {
	err = client.AppsV1().StatefulSets(namespace).Delete(context.TODO(), statefulSetName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error(errors.New("删除StatefulSet失败, " + err.Error()))
		return errors.New("删除StatefulSet失败, " + err.Error())
	}

	return nil
}

// 更新statefulset
func (s *statefulset) UpdateStatefulSet(client *kubernetes.Clientset, namespace, content string) (err error) {
	var statefulSet = &appsv1.StatefulSet{}

	err = json.Unmarshal([]byte(content), statefulSet)
	if err != nil {
		logger.Error(errors.New("反序列化失败, " + err.Error()))
		return errors.New("反序列化失败, " + err.Error())
	}

	_, err = client.AppsV1().StatefulSets(namespace).Update(context.TODO(), statefulSet, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(errors.New("更新StatefulSet失败, " + err.Error()))
		return errors.New("更新StatefulSet失败, " + err.Error())
	}
	return nil
}

// toCells 将 StatefulSet 列表转换为 DataCell 列表。
// 这个函数的存在是为了适应特定的数据处理需求，即将 Kubernetes 中的 StatefulSet 对象转换成更通用的数据结构，
// 以便在不同的上下文中使用。这种转换有助于数据的统一处理和展示。
// 参数:
//
//	std []appsv1.StatefulSet: 一个包含多个 StatefulSet 对象的切片。
//
// 返回值:
//
//	[]DataCell: 一个 DataCell 类型的切片，每个元素都是从 StatefulSet 对象转换而来。
//	            DataCell 是一种通用的数据结构，用于在不同场景下传递和展示数据。
//
// 作者: 南宫乘风
func (s *statefulset) toCells(std []appsv1.StatefulSet) []DataCell {
	// 创建一个与输入切片长度相同的 DataCell 切片，用于存储转换后的数据。
	cells := make([]DataCell, len(std))
	// 遍历输入的 StatefulSet 列表，将每个 StatefulSet 对象转换为 DataCell 对象。
	for i := range std {
		// 将当前 StatefulSet 对象转换为 DataCell 对象，并存储在相应的位置。
		cells[i] = statefulSetCell(std[i])
	}
	// 返回转换后的 DataCell 列表。
	return cells
}

// fromCells 将 DataCell 类型的切片转换为 appsv1.StatefulSet 类型的切片。
// 此函数的存在是为了解决从原始数据结构到目标 Kubernetes StatefulSet 对象的转换需求。
// 它遍历输入的 DataCell 切片，将每个 DataCell 转换为一个 StatefulSet 对象。
// 这个转换过程是必要的，因为它允许我们以一种结构化和类型安全的方式处理数据。
// 参数:
//
//	cells []DataCell - 一个 DataCell 类型的切片，每个 DataCell 都封装了一个可能的 StatefulSet 配置。
//
// 返回值:
//
//	[]appsv1.StatefulSet - 一个 appsv1.StatefulSet 类型的切片，转换自输入的 DataCell 切片。
//
// 作者: 南宫乘风
func (s *statefulset) fromCells(cells []DataCell) []appsv1.StatefulSet {
	// 创建一个与输入 cells 切片长度相同的 statefulSets 切片，用于存储转换后的 StatefulSet 对象。
	statefulSets := make([]appsv1.StatefulSet, len(cells))
	// 遍历 cells 切片，对每个 DataCell 进行转换。
	for i := range cells {
		// 将 DataCell 转换为 StatefulSet 对象，并存储在相应的索引位置。
		// 这里使用了类型断言来将 DataCell 转换为具体的 statefulSetCell 类型，然后再转换为 appsv1.StatefulSet。
		// 这种转换假设 DataCell 类型可以安全地转换为 statefulSetCell 类型，这在函数外部的上下文中应该是有效的。
		statefulSets[i] = appsv1.StatefulSet(cells[i].(statefulSetCell))
	}

	// 返回转换后的 StatefulSet 切片。
	return statefulSets
}
