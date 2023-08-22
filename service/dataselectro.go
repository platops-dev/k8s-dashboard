package service

import (
	_"fmt"
	"sort"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	nwv1     "k8s.io/api/networking/v1"
)

/*
1. 定义数据结构
*/
// dataSelect 用于封装排序、过滤、分页的数据类型
type dataSelector struct {
	GenericDataList []DataCell
	DataSelectQuery *DataSelectQuery
}

// DataCell接口，用于各种资源list的类型转换, 转换后可以使用dataselector的自定义排序方法
type DataCell interface {
	GetCreation() time.Time
	GetName() string
}

// DataSelectQuery 定义过滤和分页的属性，过滤: Name, 分页: Limit和Page
// Limit 是单页的数据条数
// Page是第几页
type DataSelectQuery struct {
	FilterQuery   *FilterQuery
	PaginateQuery *PaginateQuery
}

type FilterQuery struct {
	Name string
}

type PaginateQuery struct {
	Limit int
	Page  int
}

/*
2. 排序
实现自定义结构的排序, 需要重写Len、Swap、Less方法
*/
// Len 方法用于获取数组长度
func (d *dataSelector) Len() int {
	return len(d.GenericDataList)
}

// Swap方法用于数据中的元素在比较大小之后的位置交换，可以定义升序或者降序
func (d *dataSelector) Swap(i, j int) {
	d.GenericDataList[i], d.GenericDataList[j] = d.GenericDataList[j], d.GenericDataList[i]
}

// Less 方法用于定义数组中元素排序的“大小”的比较方式
func (d *dataSelector) Less(i, j int) bool {
	a := d.GenericDataList[i].GetCreation()
	b := d.GenericDataList[j].GetCreation()
	return b.Before(a)
}

// 重写以上三个方法后使用sort.Sort进行排序
func (d *dataSelector) Sort() *dataSelector {
	sort.Sort(d)
	return d
}

/*
2. 过滤
*/
//Filter方法用于过滤元素, 比较元素的Name属性，若包含,在返回
func (d *dataSelector) Filter() *dataSelector {
	//若Nmae的传参为空, 则返回所有元素
	//fmt.Println("这个是Filter方法看到的: ", d.DataSelectQuery.FilterQuery.Name)
	if d.DataSelectQuery.FilterQuery.Name == "" {
		// fmt.Println("我在测试dataselectro的d值", d)
		// fmt.Println("我在测试dataselectro的d值,就是上面输出的")
		return d
	} 
	// fmt.Println("没有在往下走了")
	//若Name的传参不为空, 则返回元素名中包含Name的所有元素
	FilteredList := []DataCell{}
	for _, value := range d.GenericDataList {
		matches := true
		objName := value.GetName()
		if strings.Contains(objName, d.DataSelectQuery.FilterQuery.Name) {
			matches = false
			continue
		}
		if matches {
			FilteredList = append(FilteredList, value)
		}
	}
	d.GenericDataList = FilteredList
	return d
}

/*
3. 分页
*/
// Paginate 方法用于数组分页, 根据Limit和page的传参，返回数据
func (d *dataSelector) Paginate() *dataSelector {
	limit := d.DataSelectQuery.PaginateQuery.Limit
	page := d.DataSelectQuery.PaginateQuery.Page
	//验证参数合法, 若参数不合法, 则不返回数据
	// if limit <= 0 || page <= 0 {
	// 	return d
	// }
	// 举例: 25个元素的数组, limit是10, page是3, startIndex是20, endIndex是30
	//（实际endIndex是25）
	startIndex := limit * (page - 1)
	endIndex := limit * page

	//处理最后一页, 这时候就把endIndex由30改为25了
	if len(d.GenericDataList) < endIndex {
		endIndex = len(d.GenericDataList)
	}

	d.GenericDataList = d.GenericDataList[startIndex:endIndex]
	return d
}

/*
4. 定义podCell类型, 实现DataCell接口, 用于类型转换
*/
// 定义podCell类型, 实现GetCreation 和 GetName方法后, 可进行类型转换
type podCell corev1.Pod

func (p podCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}

func (p podCell) GetName() string {
	return p.Name
}

/*
5. 定义deploymentCell类型, 实现DataCell接口, 用于类型转换
*/
type deploymentCell appsv1.Deployment

func (d deploymentCell) GetCreation() time.Time {
	return d.CreationTimestamp.Time
}

func (d deploymentCell) GetName() string {
	return d.Name
}
/*
6. 定义statefulSetCell类型, 实现DataCell接口, 用于类型转换
*/
type statefulSetCell appsv1.StatefulSet

func (s statefulSetCell) GetCreation() time.Time {
	return s.CreationTimestamp.Time
}

func (s statefulSetCell) GetName() string {
	return s.Name
}

/*
7. 定义statefulSetCell类型, 实现DataCell接口, 用于类型转换
*/
type daemonSetCell appsv1.DaemonSet

func (ds daemonSetCell) GetCreation() time.Time {
	return ds.CreationTimestamp.Time
}

func (ds daemonSetCell) GetName() string {
	return ds.Name
}

/*
8. 定义k8sNodeCell类型, 实现DataCell接口, 用于类型转换
*/
type k8sNodeCell corev1.Node

func (kn k8sNodeCell) GetCreation() time.Time {
	return kn.CreationTimestamp.Time
}

func (kn k8sNodeCell) GetName() string {
	return kn.Name
}


/*
9. 定义namespaceCell类型, 实现DataCell接口, 用于类型转换
*/
type namespaceCell corev1.Namespace

func (ns namespaceCell) GetCreation() time.Time {
	return ns.CreationTimestamp.Time
}

func (ns namespaceCell) GetName() string {
	return ns.Name
}

/*
10. 定义persistentVolumeCell类型, 实现DataCell接口, 用于类型转换
*/
type persistentVolumeCell corev1.PersistentVolume

func (pv persistentVolumeCell) GetCreation() time.Time {
	return pv.CreationTimestamp.Time
}

func (pv persistentVolumeCell) GetName() string {
	return pv.Name
}

/*
11. 定义k8sServiceCell类型, 实现DataCell接口, 用于类型转换
*/
type k8sServiceCell corev1.Service

func (svc k8sServiceCell) GetCreation() time.Time {
	return svc.CreationTimestamp.Time
}

func (svc k8sServiceCell) GetName() string {
	return svc.Name
}

/*
12. 定义k8sServiceCell类型, 实现DataCell接口, 用于类型转换
*/
type ingressCell nwv1.Ingress

func (i ingressCell) GetCreation() time.Time {
	return i.CreationTimestamp.Time
}

func (i ingressCell) GetName() string {
	return i.Name
}