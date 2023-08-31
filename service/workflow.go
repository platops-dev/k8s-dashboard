package service

import (
	"test4/dao"
	"test4/model"
	corev1 "k8s.io/api/core/v1"
)

var Workflow workflow

type workflow struct{}

// 定义WorkflowCreate结构体, 用于创建workflow需要的参数属性的定义
type WorkflowCreate struct {
	Name          string                 `json:"name"`
	Namespace     string                 `json:"namespace"`
	Replicas      int32                  `json:"replicas"`
	Image         string                 `json:"image"`
	Label         map[string]string      `json:"label"`
	Cpu           string                 `json:"cpu"`
	Memory        string                 `json:"memory"`
	Containers  []corev1.Container    	 `json:"containers"`
	ContainerPortSvc int32				 `json:"container_port_svc"`
	HealthCheck   bool                   `json:"health_check"`
	ResourceCheck bool					 `json:"resource_check"`
	HealthPath    string                 `json:"health_path"`
	Type          string                 `json:"type"`
	Port          int32                  `json:"port"`
	NodePort      int32                  `json:"node_port"`
	Hosts         map[string][]*HttpPath `json:"hosts"`
}

// 获取列表分页查询
func (wf *workflow) GetList(name string, page, limit int) (workflowResp *dao.WorkflowResp, err error) {
	workflowResp, err = dao.Workflow.GetList(name, page, limit)
	if err != nil {
		return nil, err
	}
	return workflowResp, nil
}

//获取workflow详情
func (wf *workflow) GetById(id int) (workflow *model.Workflow, err error) {
	workflow, err = dao.Workflow.GetById(id)
	if err != nil {
		return nil, err
	}
	return workflow, nil
}


//workflow名字转service名字, 添加-svc后缀
func getServiceName(workflowName string) (serviceName string) {
	return workflowName + "-svc"
}

//workflow名字转ingress名字, 添加-svc后缀
func getIngressName(workflowName string) (ingressName string) {
	return workflowName + "-ing"
}

//创建workflow
func (wf *workflow) CreateWorkflow(data *WorkflowCreate) (err error) {
	//若workflow不是ingress类型, 传入空字符串即可
	var ingressName string
	if data.Type == "Ingress" {
		ingressName = getIngressName(data.Name)
	} else {
		ingressName = ""
	}
	//组装mysql中workflow的单条数据
	workflow := &model.Workflow{
		Name: 		data.Name,
		Namespace: 	data.Namespace,
		Replicas: 	data.Replicas,
		Deployment: data.Name,
		Service: 	getServiceName(data.Name),
		Ingress: 	ingressName,
		Type: 		data.Type,
	}
	//调用dao层执行数据库添加操作
	err = dao.Workflow.Add(workflow)
	if err != nil {
		return err
	}

	//创建k8s资源
	err = createWorkflowRes(data)
	if err != nil {
		return err
	}
	return nil
}

//封装创建workflow对应的k8s资源
//小写开头的函数, 作用域只在当前包中, 不支持跨包调用
func createWorkflowRes(data *WorkflowCreate) (err error) {
	//声明service类型
	var serviceType string
	//组装deploymentCreate类型的数据
	dc := &DeployCreate{
		Name: data.Name,
		Namespace: data.Namespace,
		Replicas: data.Replicas,
		Image: data.Image,
		Label: data.Label,
		Cpu: data.Cpu,
		Memory: data.Memory,
		Containers: data.Containers,
		HealthCheck: data.HealthCheck,
		HealthPath: data.HealthPath,
		ResourceCheck: data.ResourceCheck,
	}

	//创建deployment
	err = Deployment.CreateDeployment(dc)
	if err != nil {
		return err
	}

	//判断service类型
	if data.Type != "Ingress" {
		serviceType = data.Type
	} else {
		serviceType = "ClusterIP"
	}

	//组装ServiceCreate类型的数据
	sc := &ServiceCreate{
		Name: getServiceName(data.Name),
		Namespace: data.Namespace,
		Type: serviceType,
		ContainerPort: data.ContainerPortSvc,
		Port: data.Port,
		NodePort: data.NodePort,
		Label: data.Label,
	}
	err = K8sService.CreateService(sc)
	if err != nil {
		return err
	}

	//组装IngressCreate类型的数据, 创建ingres, 只有ingress类型的workflow才有ingress资源, 所以这里做了一层判断
	if data.Type == "Ingress" {
		ic := &IngressCreate{
			Name: getIngressName(data.Name),
			Namespace: data.Namespace,
			Label: data.Label,
			Hosts: data.Hosts,
		}
		err = Ingress.CreateIngress(ic)
		if err != nil {
			return err
		}
	}
	return nil
}


//删除workflow
func (wf *workflow) DelById(id int) (err error) {
	//获取workflow资源
	workflow, err := dao.Workflow.GetById(id)
	if err != nil {
		return err
	}
	//删除k8s资源
	err = delWorkflowRes(workflow)
	if err != nil {
		return err
	}
	//删除数据库数据
	err = dao.Workflow.DelById(id)
	if err != nil {
		return err
	}
	return nil
}

//封装删除workflow对应的k8s资源
func delWorkflowRes(workflow *model.Workflow) (err error) {
	//删除deployment
	err = Deployment.DeleteDeployment(workflow.Name, workflow.Namespace)
	if err != nil {
		return err
	}
	//删除service
	err = K8sService.DeleteK8sService(getServiceName(workflow.Name), workflow.Namespace)
	if err != nil {
		return err
	}
	//删除ingress, 这里多了一层判断, 因为只有type为ingress的workflow才有ingress资源
	if workflow.Type == "Ingress" {
		err = Ingress.DeleteIngress(getIngressName(workflow.Name), workflow.Namespace)
		if err != nil {
			return err
		}
	}
	return nil
} 