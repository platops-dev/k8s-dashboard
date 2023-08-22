package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	_"strconv"
	"time"

	"github.com/wonderivan/logger"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// corev1 _"k8s.io/api/core/v1"
var Deployment deployment

type deployment struct{}

// 定义列表的返回内容, Items是deployment元素列表, Total为deployment元素数量
type DeploymentsResp struct {
	Items []appsv1.Deployment `json:"items"`
	Total int                 `json:"total"`
}

// 定义DeployCreate结构体, 用于创建deployment需要的参数属性的定义
type DeployCreate struct {
	Name          string            `json:"name"`
	Namespace     string            `json:"namespace"`
	Replicas      int32             `json:"replicas"`
	Image         string            `json:"image"`
	Label         map[string]string `json:"label"`
	Cpu           string            `json:"cpu"`
	Memory        string            `json:"memory"`
	//ContainerPort int32             `json:"container_port"`
	Containers 	  []corev1.Container `json:"containers"`
	HealthCheck   bool              `json:"health_check"`
	HealthPath    string            `json:"health_path"`
	ResourceCheck bool				`json:"resource_check"`
}

// 定义DeploysNP类型, 用于返回namespace中deployment的数量
type DeploysNp struct {
	Namespace string `json:"namespace"`
	DeployNum int    `json:"deployment_num"`
}

/*
定义DataCell 到Deployment类型转换的方法
*/
//toCells方法用于将Deployment类型数组，转换成DataCell类型数组
func (d *deployment) toCells(std []appsv1.Deployment) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = deploymentCell(std[i])
	}
	return cells
}

// fromCells 方法用于将DataCell类型数组, 转换成Deployment类型数组
/*
在 `fromCells` 方法中，我们需要将 `[]DataCell` 类型的切片转换为 `[]appsv1.Deployment` 类型的切片。但是，`DataCell` 类型是一个接口类型，它的底层类型可以是任何类型。因此，我们需要使用类型断言将其转换为 `deploymentCell` 类型，以便将其转换为 `appsv1.Deployment` 类型。
在这里，`cells[i].(deploymentCell)` 表示将 `cells[i]` 转换为 `deploymentCell` 类型，并返回其底层值。这里使用了类型断言 `(deploymentCell)`，它告诉编译器 `cells[i]` 应该是一个 `deploymentCell` 类型的值。如果 `cells[i]` 的底层类型不是 `deploymentCell`，则会在运行时触发 panic。
一旦我们将 `cells[i]` 转换为 `deploymentCell` 类型，我们就可以将其转换为 `appsv1.Deployment` 类型，并将其存储在 `deployments[i]` 中。这样做可以确保我们将正确的数据类型存储在新的切片中。
总之，这里使用了类型断言来获取底层类型并进行类型转换，以确保我们可以正确地将 `DataCell` 类型转换为 `appsv1.Deployment` 类型。
*/
func (d *deployment) fromCells(cells []DataCell) []appsv1.Deployment {
	deployments := make([]appsv1.Deployment, len(cells))
	for i := range cells {
		deployments[i] = appsv1.Deployment(cells[i].(deploymentCell))
	}
	return deployments
}

// 获取deployment 列表, 支持过滤、排序、分页
func (d *deployment) GetDeployments(filterName, namespace string, limit, page int) (deploymentsResp *DeploymentsResp, err error) {
	//获取deploymentList类型的deployment列表
	deploymentList, err := K8s.Clientset.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(errors.New("获取Deployment列表失败, " + err.Error()))
		return nil, errors.New("获取Deployment列表失败, " + err.Error())
	}
	//将deploymentList中的deploment列表(Items), 放进dataselector对象中，进行排序
	selectableData := &dataSelector{
		GenericDataList: d.toCells(deploymentList.Items),
		DataSelectQuery: &DataSelectQuery{
			FilterQuery: &FilterQuery{Name: filterName},
			PaginateQuery: &PaginateQuery{
				Limit: limit,
				Page:  page,
			},
		},
	}
	//fmt.Println(d.toCells(deploymentList.Items))

	// filtered := selectableData.Filter()
	// total := len(filtered.GenericDataList)
	// data	 := filtered.Sort().Paginate()

	//fmt.Println(selectableData)
	//fmt.Println("selectableData的值已经传了过来例如,测试这个")
	// fmt.Println( selectableData.Filter())


	//fmt.Println("我在看Filter的name值: ", selectableData.Filter().DataSelectQuery.FilterQuery.Name)
	filtered := selectableData.Filter()

	total := len(filtered.GenericDataList)

	data := filtered.Sort().Paginate()



	//将[]DataCell类型的deployment列表转为appsv1.deployment列表
	deployments := d.fromCells(data.GenericDataList)


	return &DeploymentsResp{
		Items: deployments,
		Total: total,
	}, nil
}

// 获取deployment详情
func (d *deployment) GetDeploymentDetail(deploymentName, namespace string) (deployment *appsv1.Deployment, err error) {
	deployment, err = K8s.Clientset.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		logger.Error(errors.New("获取Deployment详情失败, " + err.Error()))
		return nil, errors.New("获取Deployment详情失败, " + err.Error())
	}
	return deployment, nil
}

//设置deployment副本数
func (d *deployment) ScaleDeployment(deploymentName, namespace string, scaleNum int) (replica int32, err error) {
	//获取autoscalingv1.Scale类型的对象, 能点出当前的副本数
	scale, err := K8s.Clientset.AppsV1().Deployments(namespace).GetScale(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		logger.Error(errors.New("获取Deployment副本数信息失败, " + err.Error()))
		return 0, errors.New("获取Deployment副本数信息失败, " + err.Error())
	}
	// 修改副本数
	scale.Spec.Replicas = int32(scaleNum)
	//更新副本数，传入scale对象
	newScale, err := K8s.Clientset.AppsV1().Deployments(namespace).UpdateScale(context.TODO(), deploymentName, scale, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(errors.New("更新Deployment副本数信息失败, " + err.Error()))
		return 0, errors.New("更新Deployment副本数信息失败, " + err.Error())
	}
	return newScale.Spec.Replicas, nil
}

//创建deployment, 并接收DeployCreate对象
func (d *deployment) CreateDeployment(data *DeployCreate) (err error) {
	//将data中的数据组组装成appsv1.Deployment对象
	deployment := &appsv1.Deployment{
		// ObjectMeta 中定义资源名，命名空间以及标签
		ObjectMeta: metav1.ObjectMeta{
			Name: data.Name,
			Namespace: data.Namespace,
			Labels: data.Label,
		},
		//Spec中定义副本数、选择器、以及pod属性
		Spec: appsv1.DeploymentSpec{
			Replicas: &data.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: data.Label,
			},
			Template: corev1.PodTemplateSpec{
				//定义pod名和标签
				ObjectMeta: metav1.ObjectMeta{
					Name: data.Name,
					Labels: data.Label,
				},
				//定义容器名、镜像、和端口
				Spec: corev1.PodSpec{
					/*
					Containers: []corev1.Container{
						{
							Name: data.Name,
							Image: data.Image,
							Ports: []corev1.ContainerPort{
								{
									Name: "http",
									Protocol: corev1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
					*/
					Containers: data.Containers,
				},
			},
		},
		//Status定义资源的运行状态, 这里由于是新建，传入空的appsv1.DeploymentStatus{}对象即可
		Status: appsv1.DeploymentStatus{},
	}
	for i := range deployment.Spec.Template.Spec.Containers {
		// 判断是否打开健康检查功能，若打开，则定义ReadinessProbe和LivenessProbe
		if data.HealthCheck {
			ContainerPort := deployment.Spec.Template.Spec.Containers[i].Ports[0].ContainerPort
			// 设置第一个容器的ReadinessProbe, 因为我们pod中只有一个容器，所以直接使用index 0 即可
			// 若pod中有多个容器，则这里需要使用for循环去定义了
				deployment.Spec.Template.Spec.Containers[i].ReadinessProbe = &corev1.Probe{
					ProbeHandler: corev1.ProbeHandler{
						HTTPGet: &corev1.HTTPGetAction{
							Path: data.HealthPath,
							//intstr.IntOrString 的作用是端口可以定义为整型，也可以定义为字符串
							// Type=0 则表示该结构体实例内的数据为整型, 转json时只使用IntVal的数据
							// Type=1 则表示该结构体实例内的数据为字符串, 转json时只使用StrVal的数据
							Port: intstr.IntOrString{
								Type: 0,
								IntVal: ContainerPort,
							},
						},
					},
					//初始化等待时间
					InitialDelaySeconds: 5,
					//超时时间
					TimeoutSeconds: 5,
					//执行间隔
					PeriodSeconds: 5,
				}
				deployment.Spec.Template.Spec.Containers[i].LivenessProbe = &corev1.Probe{
					ProbeHandler: corev1.ProbeHandler{
						HTTPGet: &corev1.HTTPGetAction{
							Path: data.HealthPath,
							Port: intstr.IntOrString{
								Type: 0,
								IntVal: ContainerPort,
							},
						},
					},
					InitialDelaySeconds: 15,
					TimeoutSeconds: 5,
					PeriodSeconds: 5,
				}
		}
		if data.ResourceCheck {
			//定义容器的limit/request资源(后续可以将资源限制单独提出, 健康检查与资源限制并没有要必须一起配置)
			deployment.Spec.Template.Spec.Containers[i].Resources.Limits = map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceCPU :	resource.MustParse(data.Cpu),
				corev1.ResourceMemory :	resource.MustParse(data.Memory),
			}
			deployment.Spec.Template.Spec.Containers[i].Resources.Requests = map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceCPU : resource.MustParse(data.Cpu),
				corev1.ResourceMemory : resource.MustParse(data.Memory),
			}
		}
	}
	//调用sdk创建deployment
	_, err = K8s.Clientset.AppsV1().Deployments(data.Namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		logger.Error(errors.New("创建Deployment失败, " + err.Error()))
		return errors.New("创建Deployment失败, " + err.Error())
	}
	return nil
}

//删除deployment
func (d *deployment) DeleteDeployment(deploymentName, namespace string) (err error) {
	err = K8s.Clientset.AppsV1().Deployments(namespace).Delete(context.TODO(), deploymentName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error(errors.New("删除Deployment失败, " + err.Error()))
		return errors.New("删除Deployment失败, " + err.Error())
	}
	return nil
}


//重启deployment
func (d *deployment) RestartDeployment(deploymentName, namespace string) (err error) {
	// 此功能等同于kubectl 命令
	/*
	kubectl deploy $(service) -p \
	'{"spec":{"template":{"spec":{"containers":[{"name":"'"${service}"'","env":[{"name":"RESTART_","value":"'$(data +%s)'"}]}]}}}}'
	*/ 

	// 使用pathData Map组装数据
	/*
	老方法, 侵入性强
	patchData := map[string]interface{}{
		"spec": map[string]interface{}{
			"template": map[string]interface{}{
				"spec": map[string]interface{}{
					"containers": []map[string]interface{}{
						{
							"name": deploymentName,
								"env": []map[string]string{{
									"name": "RESTART_",
									"value": strconv.FormatInt(time.Now().Unix(), 10),
								}},
						},
					},
				},
			},
		},	
	}
	*/
	//新方法,侵入性弱
	patchData := map[string]interface{}{
		"spec": map[string]interface{}{
			"template": map[string]interface{}{
				"metadata": map[string]interface{}{
					"annotations": map[string]string{
						"kubectl.kubenetes.io/restarteAt" : metav1.Now().Format(time.RFC3339),
					},
				},
			},
		},
	}
	//序列化为字节, 因为patch方法只接收字节类型参数
	patchByte, err := json.Marshal(patchData)
	if err != nil {
		logger.Error(errors.New("JSON序列化失败, " + err.Error()))
		return errors.New("JSON序列化失败, " + err.Error())
	}
	// 调用patch方法更新deployment
	_, err = K8s.Clientset.AppsV1().Deployments(namespace).Patch(context.TODO(), deploymentName, "application/strategic-merge-patch+json", patchByte, metav1.PatchOptions{})
	if err != nil {
		logger.Error(errors.New("重启Deploment失败, " + err.Error()))
		return errors.New("重启Deployment失败, " + err.Error())
	}
	return nil
}

//更新deployment
func (d *deployment) UpdateDeployment(namespace, content string) (err error) {
	fmt.Println(content)
	var deploy = &appsv1.Deployment{}
	err = json.Unmarshal([]byte(content), deploy)
	if err != nil {
		fmt.Println(content)
		logger.Error(errors.New("反序列化失败, " + err.Error()))
		return errors.New("反序列化失败, " + err.Error())
	}
	_, err = K8s.Clientset.AppsV1().Deployments(namespace).Update(context.TODO(), deploy, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(errors.New("更新Deployment失败, " + err.Error()))
		return errors.New("更新Deployment失败, " + err.Error())
	}
	return nil
}

//获取每个namespace的deployment数量
func (d *deployment) GetDeployNumPerNP() (deploysNps []*DeploysNp, err error) {
	namespaceList, err := K8s.Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, namespace := range namespaceList.Items {
		deploymentList, err := K8s.Clientset.AppsV1().Deployments(namespace.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
		deploysNp := &DeploysNp{
			Namespace: namespace.Name,
			DeployNum: len(deploymentList.Items),
		}

		deploysNps = append(deploysNps, deploysNp)
	}
	return deploysNps, nil
}