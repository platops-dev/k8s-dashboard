package service

import (
	"context"
	"encoding/json"
	"errors"
	_"strconv"
	"time"

	"github.com/wonderivan/logger"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var StatefulSet statefulSet

type statefulSet struct{}

type StatefulSetResp struct {
	Items []appsv1.StatefulSet	`json:"items"`
	Total	int					`json:"total"`
}

type StatefulSetCreate struct {
	StatefulSetName		string				`json:"name"`
	Namespace			string				`json:"namespace"`
	Labels				map[string]string	`json:"labels"`
	VolumeMountName		string				`json:"volume_mount_name"`
	VolumeNameMountPath	string				`json:"volume_name_mount_path"`
	AccessModes			string				`json:"access_modes"`
	StorageClassName	string				`json:"storage_class_name"`
	VolumeMode			string				`json:"volume_mode"`
	Replicas			int32				`json:"replicas"`
	Image				string				`json:"image"`
	Cpu					string				`json:"cpu"`
	Memory				string				`json:"memory"`
	StorageSize			string				`json:"storage_size"`
	ContainerName		string				`json:"container_name"`
	ContainerPort		int32				`json:"container_port"`
	ContainerPortName	string				`json:"container_port_name"`
	Protocol			string				`json:"protocol"`
	HealthCheck			bool				`json:"health_check"`
	HealthPath			string				`json:"health_path"`
	HostPathCheck		bool				`json:"host_path_check"`
	ResourceCheck		bool				`json:"resource_check"`
	HostPath			string				`json:"host_path"`
}

//定义namespace 下statefulset的数量
type StatefulSetNp struct {
	Namespace			string	`json:"namespace"`
	StatefulSetNum		int		`json:"statefulset_num"`
}

//类型转换
func (s *statefulSet) tocells(std []appsv1.StatefulSet) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = statefulSetCell(std[i])
	}
	return cells
}

func (s *statefulSet) formcells(cells []DataCell) []appsv1.StatefulSet {
	statefulSet := make([]appsv1.StatefulSet, len(cells))
	for i := range cells {
		statefulSet[i] = appsv1.StatefulSet(cells[i].(statefulSetCell))
	}
	return statefulSet
}

//获取statefulSet列表，支持过滤、分页、排序
func (s *statefulSet) GetStatefulSets(filterName, namespace string, limit, page int) (statefulSetResp *StatefulSetResp, err error) {
	//获取StatefulSetList类型的statefulset列表
	StatefulSetList, err := K8s.Clientset.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(errors.New("获取Statefulset列表失败." + err.Error()))
		return nil, errors.New("获取Statefulset列表失败." + err.Error())
	}
	//将StatefulSetList中的statefulset列表(Iteams), 放进dataselector对象中，进行排序、过滤、分页
	selectableData := &dataSelector{
		GenericDataList: s.tocells(StatefulSetList.Items),
		DataSelectQuery: &DataSelectQuery{
			FilterQuery: &FilterQuery{Name: filterName},
			PaginateQuery: &PaginateQuery{
				Limit: limit,
				Page: page,
			},
		},
	}
	
	filtered := selectableData.Filter()
	total    := len(filtered.GenericDataList)
	data 	 := filtered.Sort().Paginate()

	//将[]DataCell类型的statefulset列表转换成appsv1.StatefulSet列表
	statefulSets := s.formcells(data.GenericDataList)

	return &StatefulSetResp{
		Items: statefulSets,
		Total: total,
	}, nil
}

//获取statefulset详情
func (s *statefulSet) GetStatefulSetDetail(statefulSetName, namespace string) (statefulSet *appsv1.StatefulSet, err error) {
	StatefulSet, err := K8s.Clientset.AppsV1().StatefulSets(namespace).Get(context.TODO(), statefulSetName, metav1.GetOptions{})
	if err != nil {
		logger.Error(errors.New("获取statefulset详情失败." + err.Error()))
		return nil, errors.New("获取statefulset详情失败." + err.Error())
	}
	return StatefulSet, nil
}

//设置statefulset副本数
func (s *statefulSet) ScaleStatefulSet(statefulSetName, namespace string, scalenum int) (replicas int32, err error) {
	//获取当前副本数
	scale, err := K8s.Clientset.AppsV1().StatefulSets(namespace).GetScale(context.TODO(), statefulSetName, metav1.GetOptions{})
	if err != nil {
		logger.Error(errors.New("获取副本数失败." + err.Error()))
		return 0, errors.New("获取副本数失败." + err.Error())
	}
	//设置副本数
	scale.Spec.Replicas = int32(scalenum)

	//更新副本数
	newScale, err := K8s.Clientset.AppsV1().StatefulSets(namespace).UpdateScale(context.TODO(), statefulSetName, scale, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(errors.New("更新副本数失败." + err.Error()))
		return 0, errors.New("更新副本数失败." + err.Error())
	}
	return newScale.Spec.Replicas, nil
}

//创建statefulset, 接收statefulset对象
func (s *statefulSet) CreateStatefulSet(data *StatefulSetCreate) (err error) {
	statefulSet := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: data.StatefulSetName,
			Namespace: data.Namespace,
			Labels: data.Labels,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &data.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: data.Labels,
			},
			Template: corev1.PodTemplateSpec{
				//定义pod名和标签
				ObjectMeta: metav1.ObjectMeta{
					Labels: data.Labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name: data.ContainerName,
							Image: data.Image,
							ImagePullPolicy: "IfNotPresent",
							Ports: []corev1.ContainerPort{
								{
									Name: data.ContainerPortName,
									ContainerPort: data.ContainerPort,
									Protocol: corev1.Protocol(data.Protocol),
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name: data.VolumeMountName,
									MountPath: data.VolumeNameMountPath,
								},
							},
						},
					},
				},
			},
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: data.VolumeMountName,
					},
					Spec: corev1.PersistentVolumeClaimSpec{
						AccessModes: []corev1.PersistentVolumeAccessMode{corev1.PersistentVolumeAccessMode(data.AccessModes)},
						StorageClassName: &data.StorageClassName,
						VolumeMode: (*corev1.PersistentVolumeMode)(&data.VolumeMode),
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceStorage : resource.MustParse(data.StorageSize),
							},
						},
					},
					Status: corev1.PersistentVolumeClaimStatus{},
				},
			},
		},
		Status: appsv1.StatefulSetStatus{},
	}
	//判断是否打开健康检查, 如果打开则胚子健康检查
	if data.HealthCheck {
		//设置第一个容器的ReadinessProbe, 因为我们pod中只有一个容器, 所以直接使用index0 即可
		//若pod中有多个容器, 则这里需要使用for循环去定义了
		statefulSet.Spec.Template.Spec.Containers[0].ReadinessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: data.HealthPath,
					//instr.InOrString 的作用时端口可以定义为整型, 也可以定义为字符串
					//Type=0 则表示该结构体实例内的数据为整型, 转json时只使用IntVal的数据
					//Type=1 则表示该结构体实例内的数据为字符串, 转json时只使用StrVal的数据
					Port: intstr.IntOrString{
						Type: 0,
						IntVal: data.ContainerPort,
					},
				},
			},
		}
		statefulSet.Spec.Template.Spec.Containers[0].LivenessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: data.HealthPath,
					//instr.InOrString 的作用时端口可以定义为整型, 也可以定义为字符串
					//Type=0 则表示该结构体实例内的数据为整型, 转json时只使用IntVal的数据
					//Type=1 则表示该结构体实例内的数据为字符串, 转json时只使用StrVal的数据
					Port: intstr.IntOrString{
						Type: 0,
						IntVal: data.ContainerPort,
					},
				},
			},
		}
	}
	//判断是否打开资源限额配置, 如果打开则配置
	if data.ResourceCheck {
		statefulSet.Spec.Template.Spec.Containers[0].Resources.Limits = map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceCPU : resource.MustParse(data.Cpu),
			corev1.ResourceMemory : resource.MustParse(data.Memory),
		}
		statefulSet.Spec.Template.Spec.Containers[0].Resources.Requests = map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceCPU : resource.MustParse(data.Cpu),
			corev1.ResourceMemory : resource.MustParse(data.Memory),
		}
	}
	_, err = K8s.Clientset.AppsV1().StatefulSets(data.Namespace).Create(context.TODO(), statefulSet, metav1.CreateOptions{})
	if err != nil {
		logger.Error(errors.New("创建StatefulSet失败, " + err.Error()))
		return errors.New("创建StatefulSet失败, " + err.Error())
	}
	return nil	
}

//删除statefulset
func (s *statefulSet) DeleteStatefulSet(statefulSetName, namespace string) (err error) {
	err = K8s.Clientset.AppsV1().StatefulSets(namespace).Delete(context.TODO(), statefulSetName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error(errors.New("删除statefulset失败." + err.Error()))
		return errors.New("删除statefulset失败." + err.Error())
	}
	return nil
}

//重启statefulset
func (s *statefulSet) RestartStatefulSet(statefulSetName, namespace string) (err error) {
	//此功能等同于kubectl命令
	// kubectp statefulset ${service} -p \
	//'{"spec":{"template":{"spec":{"containers":[{"name":"'"${service}"'","env":[{"name":"RESTART_","value":"'$(data +%s)'"}]}]}}}}'

	//使用patchData Map组装数据
	/*原来重启方式 
	patchData := map[string]interface{}{
		"spec": map[string]interface{}{
			"template": map[string]interface{}{
				"spec": map[string]interface{}{
					"containers": []map[string]interface{}{
						{
							"name": statefulSetName,
							"env": []map[string]string{
								{
									"name": "RESTART_",
									"value": strconv.FormatInt(time.Now().Unix(), 10),
								},
							},
						},
					},
				},
			},
		},
	}
	*/
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
	
	//序列化为字节， 因为patch方法只接收字节类型参数
	patchByte, err := json.Marshal(patchData)
	if err != nil {
		logger.Error(errors.New("json序列化失败." + err.Error()))
		return errors.New("json序列化失败." + err.Error())
	}
	//调用patch方法更新statefulset
	_, err = K8s.Clientset.AppsV1().StatefulSets(namespace).Patch(context.TODO(), statefulSetName, "application/strategic-merge-patch+json", patchByte, metav1.PatchOptions{})
	if err != nil {
		logger.Error(errors.New("重启statefulset失败." + err.Error()))
		return errors.New("重启statefulset失败." + err.Error())
	}
	return nil
}

//更新statefulset
func (s *statefulSet) UpdateStatefulSet(namespace, content string) (err error) {
	var   statefulSet = &appsv1.StatefulSet{}

	err = json.Unmarshal([]byte(content), statefulSet)
	if err != nil {
		logger.Error(errors.New("json反序列化失败." + err.Error()))
		return errors.New("json反序列化失败." + err.Error())
	}

	_, err = K8s.Clientset.AppsV1().StatefulSets(namespace).Update(context.TODO(), statefulSet, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(errors.New("更新statefulset失败." + err.Error()))
		return errors.New("更新statefulset失败." + err.Error())
	}
	return nil
}	

func (s *statefulSet) GetStatefulSetsNumPerNp() (statefulSetNps []*StatefulSetNp, err error) {
	namespaceList, err := K8s.Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, namespace := range namespaceList.Items {
		statefulSetList, err := K8s.Clientset.AppsV1().StatefulSets(namespace.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
		statefulSetNp := StatefulSetNp{
			Namespace: namespace.Name,
			StatefulSetNum: len(statefulSetList.Items),
		}
		statefulSetNps = append(statefulSetNps, &statefulSetNp)
	}
	return statefulSetNps, nil
}