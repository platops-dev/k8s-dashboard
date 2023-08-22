package service

import (
	"context"
	"errors"

	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var PersistentVolume persistentVolume

type persistentVolume struct{}

type PersistentVolumeResp struct {
	Items []corev1.PersistentVolume	`json:"items"`
	Total	int						`json:"total"`
}

//类型转换
func (pv *persistentVolume) toCells(std []corev1.PersistentVolume) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = persistentVolumeCell(std[i])
	}
	return cells
}

func (pv *persistentVolume) fromCells(cells []DataCell) []corev1.PersistentVolume {
	persistentVolume := make([]corev1.PersistentVolume, len(cells))
	for i := range cells {
		persistentVolume[i] = corev1.PersistentVolume(cells[i].(persistentVolumeCell))
	}
	return persistentVolume
}

func (pv *persistentVolume) GetPersistentVolumes(filterName string, limit, page int) (persistentVolumeResp *PersistentVolumeResp, err error) {
	PersistentVolumeList, err := K8s.Clientset.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(errors.New("获取PersistentVolumeList 列表失败." + err.Error()))
		return nil, errors.New("获取PersistentVolumeList 列表失败." + err.Error()) 
	}

	selectableData := &dataSelector{
		GenericDataList: pv.toCells(PersistentVolumeList.Items),
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

	persistentVolumeResps := pv.fromCells(data.GenericDataList)

	return &PersistentVolumeResp{
		Items: persistentVolumeResps,
		Total: total,
	}, nil
}

func (pv *persistentVolume) GetPersistentVolumeDetail(persistentVolumeName string) (persistentVolume *corev1.PersistentVolume, err error) {
	persistentVolume, err = K8s.Clientset.CoreV1().PersistentVolumes().Get(context.TODO(), persistentVolumeName, metav1.GetOptions{})
	if err != nil {
		logger.Error(errors.New("获取PersistentVolume: %s 详情失败." + err.Error()), persistentVolumeName)
		return nil, errors.New("获取PersistentVolume  详情失败." + err.Error()) 
	}
	return persistentVolume, nil
}

func (pv *persistentVolume) DeletePersistentVolume(persistentVolumeName string) (err error) {
	err = K8s.Clientset.CoreV1().PersistentVolumes().Delete(context.TODO(), persistentVolumeName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error(errors.New("删除PersistentVolume: %s 失败." + err.Error()), persistentVolumeName)
		return errors.New("获取PersistentVolume 失败." + err.Error()) 
	}
	return nil
}