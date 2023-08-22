package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/wonderivan/logger"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var DaemonSet daemonSet

type daemonSet struct{}

type DaemonSetResp struct {
	Items []appsv1.DaemonSet	`json:"items"`
	Total	int					`json:"total"`
}

type DaemonSetCreate struct {
	DaemonSetName	string				`json:"daemonset_name"`
	Namespace		string				`json:"namespace"`
	Labels			map[string]string	`json:"labels"`
	Containers		[]corev1.Container	`json:"containers"`
	Volume			[]corev1.Volume		`json:"volume"`
	Limit           string				`json:"limit"`
	Page			string				`json:"page"`
	Health_Check	bool				`json:"health_check"`
	Health_Path		string				`json:"health_path"`
	Resource_Check	bool				`json:"resource_check"`
	Cpu				string				`json:"cpu"`
	Memory			string				`json:"memory"`
}


type DaemonSetsNp struct {
	Namespace		string	`json:"namespace"`
	DaemonSetNum	int		`json:"daemonset_num"`
}

//类型转换
func (ds *daemonSet) tocells(std []appsv1.DaemonSet) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = daemonSetCell(std[i])
	}
	return cells
}

func (ds *daemonSet) formcells(cells []DataCell) []appsv1.DaemonSet {
	daemonSet := make([]appsv1.DaemonSet, len(cells))
	for i := range cells {
		daemonSet[i] = appsv1.DaemonSet(cells[i].(daemonSetCell))
	}
	return daemonSet
}


//获取Daemonset列表，支持过滤、排序、分页
func (ds *daemonSet) GetDaemonSets(filterName, namespace string, limit, page int) (daemonSetResp *DaemonSetResp, err error ) {
	DaemonSetList, err := K8s.Clientset.AppsV1().DaemonSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(errors.New("获取daemonset 列表失败." + err.Error()))
		return nil, errors.New("获取daemonset 列表失败." + err.Error())
	}
	//将获取到的DaemonSetList中的daemonset列表(Items)，放进dataselector对象中，进行排序、过滤、分页
	selectableData := &dataSelector{
		GenericDataList: ds.tocells(DaemonSetList.Items),
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
	data     := filtered.Sort().Paginate()

	daemonSetResps := ds.formcells(data.GenericDataList)

	return &DaemonSetResp{
		Items: daemonSetResps,
		Total: total,
	}, nil
}

//获取Daemonset详情
func (ds *daemonSet) GetDaemonSetDetail(daemonSetName, namespace string) (daemonSet *appsv1.DaemonSet, err error) {
	DaemonSet, err := K8s.Clientset.AppsV1().DaemonSets(namespace).Get(context.TODO(), daemonSetName, metav1.GetOptions{})
	if err != nil {
		logger.Error(errors.New("获取daemonset 详情失败." + err.Error()))
		return nil, errors.New("获取daemonset 详情失败." + err.Error())
	}
	return DaemonSet, nil
}

//创建Daemonset, 接收DaemonCreate对象
func (ds *daemonSet) CreateDaemonSet(data *DaemonSetCreate) (err error) {
	daemonSet := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: data.DaemonSetName,
			Namespace: data.Namespace,
			Labels: data.Labels,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: data.Labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: data.Labels,
				},
				Spec: corev1.PodSpec{
					Containers: data.Containers,
					Volumes: data.Volume,
				},
				
			},
		},
		Status: appsv1.DaemonSetStatus{},
	}

	for i := range daemonSet.Spec.Template.Spec.Containers {
		if data.Health_Check {
			Container_Port := daemonSet.Spec.Template.Spec.Containers[i].Ports[0].ContainerPort
			daemonSet.Spec.Template.Spec.Containers[i].ReadinessProbe = &corev1.Probe{
				ProbeHandler: corev1.ProbeHandler{
					HTTPGet: &corev1.HTTPGetAction{
						Path: data.Health_Path,
						//instr.InOrString 的作用是端口可以定义为整型, 也可以定义为字符串
						//Type=0 则表示该结构体实例内的数据为整型, 转json时只使用IntVal的数据
						//Type=1 则表示该结构体实例内的数据为字符串, 转json时只使用StrVal的数据
						Port: intstr.IntOrString{
							Type: 0,
							IntVal: Container_Port,
						},
					},
				},
			}
			daemonSet.Spec.Template.Spec.Containers[i].LivenessProbe = &corev1.Probe{
				ProbeHandler: corev1.ProbeHandler{
					HTTPGet: &corev1.HTTPGetAction{
						Path: data.Health_Path,
						//instr.InOrString 的作用是端口可以定义为整型, 也可以定义为字符串
						//Type=0 则表示该结构体实例内的数据为整型, 转json时只使用IntVal的数据
						//Type=1 则表示该结构体实例内的数据为字符串, 转json时只使用StrVal的数据
						Port: intstr.IntOrString{
							Type: 0,
							IntVal: Container_Port,
						},
					},
				},
			}
		}
		if data.Resource_Check {
			daemonSet.Spec.Template.Spec.Containers[i].Resources.Limits = map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceCPU : resource.MustParse(data.Cpu),
				corev1.ResourceMemory : resource.MustParse(data.Memory),
			}
			daemonSet.Spec.Template.Spec.Containers[i].Resources.Requests = map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceCPU : resource.MustParse(data.Cpu),
				corev1.ResourceMemory : resource.MustParse(data.Memory),
			}
		}
	}
		
	
	_, err = K8s.Clientset.AppsV1().DaemonSets(data.Namespace).Create(context.TODO(), daemonSet, metav1.CreateOptions{})
	if err != nil {
		logger.Error(errors.New("创建DaemonSet失败." + err.Error()))
		return errors.New("创建DaemonSet失败." + err.Error())
	}
	return nil
}

//删除daemonset
func (ds *daemonSet) DeleteDaemonSet(daemonSetName, namespace string) (err error) {
	err = K8s.Clientset.AppsV1().DaemonSets(namespace).Delete(context.TODO(), daemonSetName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error(errors.New("删除daemonset 详情失败." + err.Error()))
		return errors.New("删除daemonset 详情失败." + err.Error())
	}
	return nil
}


//重启daemonset
func (ds *daemonSet) RestartDaemonSet(daemonSetName, namespace string) (err error) {
	patchData := map[string]interface{}{
		"spec": map[string]interface{}{
			"template": map[string]interface{}{
				"metadata": map[string]interface{}{
					"annotations": map[string]string{
						"kubectl.kubenetes.io/restarteAt": metav1.Now().Format(time.RFC3339),
					},
				},
			},
		},
	}
	patchByte, err := json.Marshal(patchData)
	if err != nil {
		logger.Error(errors.New("序列化json失败." + err.Error()))
		return errors.New("序列化json失败." + err.Error())
	}
	_, err = K8s.Clientset.AppsV1().DaemonSets(namespace).Patch(context.TODO(), daemonSetName, "application/strategic-merge-patch+json", patchByte, metav1.PatchOptions{})
	if err != nil {
		logger.Error(errors.New("重启daemonset 详情失败." + err.Error()))
		return errors.New("重启daemonset 详情失败." + err.Error())
	}
	return nil
}

//更新daemonset
func (ds *daemonSet) UpdateDaemonSet(namespace, content string) (err error) {
	var daemonSet = &appsv1.DaemonSet{}

	err = json.Unmarshal([]byte(content), daemonSet)
	if err != nil {
		logger.Error(errors.New("反序列化json失败." + err.Error()))
		return errors.New("反序列化json失败." + err.Error())
	}

	_, err = K8s.Clientset.AppsV1().DaemonSets(namespace).Update(context.TODO(), daemonSet, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(errors.New("更新daemonset 失败." + err.Error()))
		return errors.New("更新daemonset 失败." + err.Error())
	}
	return nil
}

//获取每个namespace的DaemonSet数量
func (ds *daemonSet) GetDaemonSetNumPerNp() (daemonSetsNps []*DaemonSetsNp, err error) {
	namespaceList, err := K8s.Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, namespace := range namespaceList.Items {
		DaemonSetList, err := K8s.Clientset.AppsV1().DaemonSets(namespace.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return nil, err
		}

		daemonSetsNp := &DaemonSetsNp{
			Namespace: namespace.Name,
			DaemonSetNum: len(DaemonSetList.Items),
		}

		daemonSetsNps = append(daemonSetsNps, daemonSetsNp)
	}
	return daemonSetsNps, nil
}