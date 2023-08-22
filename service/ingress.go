package service

import (
	"context"
	"encoding/json"
	"errors"
	_"fmt"

	"github.com/wonderivan/logger"
	nwv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Ingress ingress

type ingress struct{}

type IngressResp struct {
	Items	[]nwv1.Ingress	`json:"items"`
	Total	int				`json:"total"`
}

//定义IngressCreate结构体, 用于创建ingress需要的参数属性的定义
type IngressCreate struct {
	Name		string	`json:"name"`
	Namespace	string	`json:"namespace"`
	Label		map[string]string	`json:"label"`
	Hosts		map[string][]*HttpPath	`json:"hosts"`
}

type HttpPath struct {
	Path	string	`json:"path"`
	PathType	nwv1.PathType	`json:"path_type"`
	ServiceName	string	`json:"service_name"`
	ServicePort	int32	`json:"service_port"`	
}



//数据类型转换
func (i *ingress) toCells(std []nwv1.Ingress) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = ingressCell(std[i])
	}
	return cells
}

func (i *ingress) fromCells(cells []DataCell) []nwv1.Ingress {
	ingress := make([]nwv1.Ingress, len(cells))
	for i := range cells {
		ingress[i] = nwv1.Ingress(cells[i].(ingressCell))
	}
	return ingress
}

func (i *ingress) GetIngress(filterName, namespace string, limit, page int) (ingressResp *IngressResp, err error) {
	IngressList, err := K8s.Clientset.NetworkingV1().Ingresses(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(errors.New("获取IngressList列表失败, " + err.Error()))
		return nil, errors.New("获取IngressList列表失败, " + err.Error())
	}

	//将获取到的IngressList中的ingress列表(Items)，放入dataselector对象中进行排序、过滤、分页
	selectableData := &dataSelector{
		GenericDataList: i.toCells(IngressList.Items),
		DataSelectQuery: &DataSelectQuery{
			FilterQuery: &FilterQuery{Name: filterName},
			PaginateQuery: &PaginateQuery{
				Limit: limit,
				Page: page,
			},
		},
	}

	filtered := selectableData.Filter()
	total := len(filtered.GenericDataList)
	data := filtered.Sort().Paginate()

	ingressResps := i.fromCells(data.GenericDataList)

	return &IngressResp{
		Items: ingressResps,
		Total: total,
	}, nil
}

func (i *ingress) GetIngressDetail(ingressName, namespace string) (ingress *nwv1.Ingress, err error) {
	Ingress, err := K8s.Clientset.NetworkingV1().Ingresses(namespace).Get(context.TODO(), ingressName, metav1.GetOptions{})
	if err != nil {
		logger.Error(errors.New("获取Ingress 详情失败, " + err.Error()))
		return nil, errors.New("获取Ingress 详情失败, " + err.Error())
	}
	return Ingress, nil
}

func (i *ingress) DeleteIngress(ingressName, namespace string) (err error) {
	err = K8s.Clientset.NetworkingV1().Ingresses(namespace).Delete(context.TODO(), ingressName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error(errors.New("删除Ingress: %s 失败, " + err.Error()), ingressName)
		return errors.New("删除 Ingress 失败, " + err.Error())
	}
	return nil
}

func (i *ingress) CreateIngress(data *IngressCreate) (err error) {
	//声明nwv1.IngressRule 和 nwv1.HTTPIngressPath变量, 后面组装数据中用到
	var ingressRules []nwv1.IngressRule
	var httpIngressPaths []nwv1.HTTPIngressPath

	//将data中的数据组装成nwv1.Ingress对象
	ingress := &nwv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name: data.Name,
			Namespace: data.Namespace,
			Labels: data.Label,
		},
		Status: nwv1.IngressStatus{},
	}
	//第一层for循环是将host组转成nwv1.IngressRule类型的对象
	// 一个host对应一个ingressRule， 每个ingressRule中包含一个host和多个path
	for key, value := range data.Hosts {
		ir := nwv1.IngressRule{
			Host: key,
			//这里将nwv1.HTTPIngressRuleValue类型中的Paths置为空, 后面组装好数据在赋值
			IngressRuleValue: nwv1.IngressRuleValue{
				HTTP: &nwv1.HTTPIngressRuleValue{Paths: nil},
			},
		}
	// 	//第二层for循环是将path组装成nwv1.HTTPIngressPath类型的对象
		for _, httpPath := range value {
			hip := nwv1.HTTPIngressPath{
				Path: httpPath.Path,
				PathType: &httpPath.PathType,
				Backend: nwv1.IngressBackend{
					Service: &nwv1.IngressServiceBackend{
						Name: httpPath.ServiceName,
						Port: nwv1.ServiceBackendPort{
							Number: httpPath.ServicePort,
						},
					},
				},
			}

			//将每个hip对象组装成数组
			httpIngressPaths = append(httpIngressPaths, hip)
			
	 	}

		//给Paths赋值， 因为前面置为空了
		ir.IngressRuleValue.HTTP.Paths = httpIngressPaths
		//将每个ir对象组装成数组, 这个ir对象就是IngressRule, 每个元素是一个host和多个path
		ingressRules = append(ingressRules, ir)
	}
	//将ingressRules对象加入到ingress的规则中
	ingress.Spec.Rules = ingressRules

	//fmt.Println(ingress.Spec.Rules)
	//创建ingress
	_, err = K8s.Clientset.NetworkingV1().Ingresses(data.Namespace).Create(context.TODO(), ingress, metav1.CreateOptions{})
	if err != nil {
		logger.Error(errors.New("创建Namespace: %s 下的Ingress: %s 失败, " + err.Error()), data.Namespace, data.Name)
		return errors.New("创建 Ingress 失败, " + err.Error())
	}
	return nil
}

func (i *ingress) UpdateIngress(namespace, content string) (err error)  {
	var ingress = &nwv1.Ingress{}

	err = json.Unmarshal([]byte(content), ingress)
	if err != nil {
		logger.Error(errors.New("JONS反序列化失败." + err.Error()))
		return errors.New("JONS反序列化失败." + err.Error())
	}
	_, err = K8s.Clientset.NetworkingV1().Ingresses(namespace).Update(context.TODO(), ingress, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(errors.New("更新Namespace: %s 下的Ingress: %s 失败, " + err.Error()), namespace, ingress.Name)
		return errors.New("更新 Ingress 失败, " + err.Error())
	}
	return nil
}