package service

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"k8s.io/apimachinery/pkg/api/resource"

	"k8s.io/apimachinery/pkg/util/intstr"

	corev1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/aryming/logger"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/kubernetes"
)

/**
 * @Author: 南宫乘风
 * @Description:
 * @File:  deployment.go
 * @Email: 1794748404@qq.com
 * @Date: 2025-03-27 15:45
 */

var Deployment deployment

type deployment struct{}

// DeploymentsResp 定义列表的返回内容，Items是deployment元素列表，Total为deployment元素数量
type DeploymentsResp struct {
	Items []appsv1.Deployment `json:"items"`
	Total int                 `json:"total"`
}

// DeployCreate 定义DeployCreate结构体，用于创建deployment需要的参数属性的定义
type DeployCreate struct {
	Name          string            `json:"name"`
	Namespace     string            `json:"namespace"`
	Replicas      int32             `json:"replicas"`
	Image         string            `json:"image"`
	Labels        map[string]string `json:"labels"`
	Cpu           string            `json:"cpu"`
	Memory        string            `json:"memory"`
	ContainerPort int32             `json:"container_port"`
	HealthCheck   bool              `json:"health_check"`
	HealthPath    string            `json:"health_path"`
	Cluster       string            `json:"cluster"`
}

// DeploysNp 定义DeploysNp类型，用于返回namespace中deployment的数量
type DeploysNp struct {
	Namespace string `json:"namespace"`
	DeployNum int    `json:"deployment_num"`
}

func (d *deployment) GetDeployments(client *kubernetes.Clientset, filterName string, namespace string, limit int, page int) (*DeploymentsResp, error) {
	// 获取deployment列表
	deploymentList, err := client.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(errors.New("获取Deployment列表失败, " + err.Error()))
		return nil, errors.New("获取Deployment列表失败, " + err.Error())
	}
	//将deploymentList中的deployment列表(Items)，放进dataselector对象中，进行排序
	selectableData := &dataSelector{
		GenericDateSelect: d.toCells(deploymentList.Items),
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

	//将[]DataCell类型的deployment列表转为appsv1.deployment列表
	deployments := d.fromCells(data.GenericDateSelect)

	return &DeploymentsResp{
		Items: deployments,
		Total: total,
	}, nil
}

// ScaleDeployment 设置deployment副本数
func (d *deployment) ScaleDeployment(client *kubernetes.Clientset, deploymentName, namespace string, scaleNum int) (replicas int32, err error) {
	// 获取autoscalingV1接口的对象，能点出当前的副本数
	scale, err := client.AppsV1().Deployments(namespace).GetScale(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		logger.Error(errors.New("获取deployment副本数失败, " + err.Error()))
		return 0, errors.New("获取deployment副本数失败, " + err.Error())
	}
	// 修改deployment副本数
	scale.Spec.Replicas = int32(scaleNum)
	// 更新deployment副本数，传入scale对象
	newScale, err := client.AppsV1().Deployments(namespace).UpdateScale(context.TODO(), deploymentName, scale, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(errors.New("更新deployment副本数失败, " + err.Error()))
		return 0, errors.New("更新deployment副本数失败, " + err.Error())
	}
	return newScale.Spec.Replicas, nil
}

// CreateDeployment 创建deployment,接收DeployCreate对象
func (d *deployment) CreateDeployment(client *kubernetes.Clientset, deployCreate *DeployCreate) (err error) {
	//	将data中的属性组装成appsv1.Deployment对象
	deployment := &appsv1.Deployment{
		//ObjectMeta中定义资源名、命名空间以及标签
		ObjectMeta: metav1.ObjectMeta{
			Name:      deployCreate.Name,
			Namespace: deployCreate.Namespace,
			Labels:    deployCreate.Labels,
		}, //Spec中定义副本数、Pod模板
		Spec: appsv1.DeploymentSpec{
			Replicas: &deployCreate.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: deployCreate.Labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   deployCreate.Name,
					Labels: deployCreate.Labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  deployCreate.Name,
							Image: deployCreate.Image,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",             //容器端口名称
									Protocol:      corev1.ProtocolTCP, //协议
									ContainerPort: 80,                 //容器端口
								},
							},
						},
					},
				},
			},
		},
		//Status定义资源的运行状态，这里由于是新建，传入空的appsv1.DeploymentStatus{}对象即可
		Status: appsv1.DeploymentStatus{},
	}

	//判断是否打开健康检查功能，若打开，则定义ReadinessProbe和LivenessProbe
	if deployCreate.HealthCheck { //如果打开健康检查
		//设置第一个容器的ReadinessProbe，因为我们pod中只有一个容器，所以直接使用index 0即可
		//若pod中有多个容器，则这里需要使用for循环去定义了

		deployment.Spec.Template.Spec.Containers[0].ReadinessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: deployCreate.HealthPath,
					Port: intstr.IntOrString{
						Type:   0,
						IntVal: deployCreate.ContainerPort,
					},
				},
			},
			//"#初始化等待时间
			InitialDelaySeconds: 5,
			//超时时间
			TimeoutSeconds: 5,
			//执行间隔
			PeriodSeconds: 5,
		}
		deployment.Spec.Template.Spec.Containers[0].LivenessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: deployCreate.HealthPath,
					Port: intstr.IntOrString{
						Type:   0,
						IntVal: deployCreate.ContainerPort,
					},
				},
			},
			InitialDelaySeconds: 15,
			TimeoutSeconds:      5,
			PeriodSeconds:       5,
		}
		//定义容器的limit和request资源
		deployment.Spec.Template.Spec.Containers[0].Resources.Limits = map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceCPU:    resource.MustParse(deployCreate.Cpu),
			corev1.ResourceMemory: resource.MustParse(deployCreate.Memory),
		}
		deployment.Spec.Template.Spec.Containers[0].Resources.Requests = map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceCPU:    resource.MustParse(deployCreate.Cpu),
			corev1.ResourceMemory: resource.MustParse(deployCreate.Memory),
		}
	}
	// 调用sdk创建deployment
	_, err = client.AppsV1().Deployments(deployment.Namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		logger.Error(errors.New("创建deployment失败, " + err.Error()))
		return errors.New("创建deployment失败, " + err.Error())
	}
	return nil
}

// RestartDeployment 重启deployment
func (d *deployment) RestartDeployment(client *kubernetes.Clientset, deploymentName, namespace string) (err error) {
	// 此功能等同于一下kubectl命令
	//使用patchData Map组装数据
	patchData := map[string]interface{}{
		"spec": map[string]interface{}{
			"template": map[string]interface{}{
				"spec": map[string]interface{}{
					"containers": []map[string]interface{}{
						{
							"name": deploymentName,
							"env": []map[string]string{
								{
									"name":  "RESTART_",
									"value": strconv.FormatInt(time.Now().Unix(), 10),
								},
							},
						},
					},
				},
			},
		},
	}
	//序列化为字节，因为patch方法只接收字节类型参数
	patchBytes, err := json.Marshal(patchData)

	if err != nil {
		logger.Error(errors.New("序列化patchData失败, " + err.Error()))
		return errors.New("序列化patchData失败, " + err.Error())
	}
	//调用patch方法更新deployment
	_, err = client.AppsV1().Deployments(namespace).Patch(context.TODO(), deploymentName, "application/strategic-merge-patch+json", patchBytes, metav1.PatchOptions{})
	if err != nil {
		logger.Error(errors.New("重启deployment失败, " + err.Error()))
		return errors.New("重启deployment失败, " + err.Error())
	}
	return nil
}

// GetDeploymentDetail 获取deployment详情
func (d *deployment) GetDeploymentDetail(client *kubernetes.Clientset, namespace string, name string) (deployment *appsv1.Deployment, err error) {
	deployment, err = client.AppsV1().Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logger.Error(errors.New("获取deployment详情失败, " + err.Error()))
		return nil, errors.New("获取deployment详情失败, " + err.Error())
	}
	return deployment, nil
}

func (d *deployment) DeleteDeployment(client *kubernetes.Clientset, name string, namespace string) (err error) {
	err = client.AppsV1().Deployments(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		logger.Error(errors.New("删除deployment失败, " + err.Error()))
		return errors.New("删除deployment失败, " + err.Error())
	}
	return nil
}

// GetDeployNumPerNp 获取每个namespace的deployment数量
func (d *deployment) GetDeployNumPerNp(client *kubernetes.Clientset) (deploysNps []*DeploysNp, err error) {
	namespaceList, err := client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(errors.New("获取namespace列表失败, " + err.Error()))
		return nil, errors.New("获取namespace列表失败, " + err.Error())
	}
	for _, namespace := range namespaceList.Items {
		deploymentList, err := client.AppsV1().Deployments(namespace.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			logger.Error(errors.New("获取namespace:" + namespace.Name + "的deployment列表失败, " + err.Error()))
			return nil, errors.New("获取namespace:" + namespace.Name + "的deployment列表失败, " + err.Error())
		}
		deploysNp := &DeploysNp{
			Namespace: namespace.Name,
			DeployNum: len(deploymentList.Items),
		}

		deploysNps = append(deploysNps, deploysNp)
	}
	return deploysNps, nil
}

func (d *deployment) toCells(std []appsv1.Deployment) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = deploymentCell(std[i])
	}
	return cells
}

func (d *deployment) fromCells(cells []DataCell) []appsv1.Deployment {
	std := make([]appsv1.Deployment, len(cells))
	for i := range std {
		std[i] = appsv1.Deployment(cells[i].(deploymentCell))
	}
	return std
}

// UpdateDeployment 更新deployment
func (d *deployment) UpdateDeployment(client *kubernetes.Clientset, namespace, content string) (err error) {
	var deploy = &appsv1.Deployment{}

	err = json.Unmarshal([]byte(content), deploy)
	if err != nil {
		logger.Error(errors.New("反序列化失败, " + err.Error()))
		return errors.New("反序列化失败, " + err.Error())
	}

	_, err = client.AppsV1().Deployments(namespace).Update(context.TODO(), deploy, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(errors.New("更新Deployment失败, " + err.Error()))
		return errors.New("更新Deployment失败, " + err.Error())
	}
	return nil
}
