package service

import (
	"sort"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
)

/**
 * @Author: 南宫乘风
 * @Description:定义数据结构
 * @File:  dataselect.go
 * @Email: 1794748404@qq.com
 * @Date: 2025-03-20 14:38
 */

// dataSelect 用于封装排序、过滤、分页的数据类型
type dataSelector struct {
	GenericDateSelect []DataCell       // 接口
	dataSelectQuery   *DataSelectQuery // 结构体
}

// DataCell 用于各种资源list的类型转换，转换后可以使用dataSelector的自定义排序方法
type DataCell interface {
	GetCreation() time.Time
	GetName() string
}

// DataSelectQuery 定义过滤和分页的属性，过滤：Name， 分页：Limit和Page
// Limit是单页的数据条数
// Page是第几页
type DataSelectQuery struct {
	FilterQuery     *FilterQuery
	PaginationQuery *PaginationQuery
}

type FilterQuery struct {
	Name string
}

type PaginationQuery struct {
	Limit int
	Page  int
}

//实现自定义结构的排序，需要重写Len、Swap、Less方法

// Len 方法用于获取数据长度
func (d *dataSelector) Len() int {
	return len(d.GenericDateSelect)
}

// Swap 方法用于数组中的元素在比较大小后的位置交换，可定义升序或降序   i j 是切片的下标
func (d *dataSelector) Swap(i, j int) {
	// 交换GenericDateSelect数组中的第i个和第j个元素
	d.GenericDateSelect[i], d.GenericDateSelect[j] = d.GenericDateSelect[j], d.GenericDateSelect[i]
}

// Less 方法用于定义数组中元素排序的“大小”的比较方式
// Less 方法返回true表示第i个元素小于第j个元素，返回false表示第i个元素大于第j个元素
func (d *dataSelector) Less(i, j int) bool {
	a := d.GenericDateSelect[i].GetCreation()
	b := d.GenericDateSelect[j].GetCreation()
	return b.Before(a)
}

// Sort 重新以上3个方法，使用sort.Sort()方法进行排序
func (d *dataSelector) Sort() *dataSelector {
	// 使用sort.Sort()方法进行排序
	sort.Sort(d)
	return d
}

// 过滤

// Filter 方法用于过滤元素，比较元素的Name属性，若包含，再返回
func (d *dataSelector) Filter() *dataSelector {
	//如Name的传参为空，则返回所有元素
	if d.dataSelectQuery.FilterQuery.Name == "" {
		return d
	}
	// 若Name的传参不为空，则返回元素中包含Name的元素
	var filteredList []DataCell
	for _, item := range d.GenericDateSelect {
		matched := true
		objName := item.GetName()
		if !strings.Contains(objName, d.dataSelectQuery.FilterQuery.Name) {
			matched = false
			continue
		}
		if matched {
			filteredList = append(filteredList, item)
		}
	}
	d.GenericDateSelect = filteredList // 返回过滤后的元素
	return d
}

// 分页

// Paginate 方法用于数组分页，根据Limit和Page的传参，返回数据
func (d *dataSelector) Paginate() *dataSelector {
	limit := d.dataSelectQuery.PaginationQuery.Limit
	page := d.dataSelectQuery.PaginationQuery.Page
	// 验证参数合法，若不合法，则返回所有元素
	if limit < 1 || page < 1 {
		return d
	}
	// 举例：25个元素的数组，limit是10，page是3，startIndex是20，endIndex是30（实际上endIndex是25）、
	startIndex := (page - 1) * limit
	endIndex := page * limit

	// 处理最后一页，这时候就把endIndex由30改为25了
	if len(d.GenericDateSelect) < endIndex {
		endIndex = len(d.GenericDateSelect)
	}
	d.GenericDateSelect = d.GenericDateSelect[startIndex:endIndex]
	return d
}

// 定义podCell 类型，实现两个方法GetCreation和GetName，可进行类型转换
type podCell corev1.Pod

func (p podCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}

func (p podCell) GetName() string {
	return p.Name
}
