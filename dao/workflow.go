package dao

import (
	"errors"
	"test4/db"
	"test4/model"

	"github.com/wonderivan/logger"
)

var Workflow workflow

type workflow struct{}

// 定义列表返回内容, Items是workflow元素列表, Total为workflow元素数量
type WorkflowResp struct {
	Items []*model.Workflow	`json:"items"`
	Total	int				`json:"total"`
}

//获取列表分页查询
func (wf *workflow) GetList(name string, page, limit int) (workflowResp *WorkflowResp, err error) {
	//定义分页数据的起始位置
	startSet := (page -1) * limit

	//定义数据库查询返回内容
	var workflowList []*model.Workflow

	//数据库查询, limit方法用于限制条数, offset方法设置起始位置
	tx := db.GORM.Where("name like ?", "%" + name + "%").
	Limit(limit).
	Offset(startSet).
	Order("id desc").
	Find(&workflowList)

	//gorm会默认把空数据也放到err中, 故这里要排除空数据的情况
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		logger.Error("获取workflow列表失败," + tx.Error.Error())
		return nil, errors.New("获取workflow列表失败," + tx.Error.Error())
	}

	return &WorkflowResp{
		Items: workflowList,
		Total: len(workflowList),
	}, nil
}

//获取workflow单条数据
func (wf *workflow) GetById(id int) (workflow *model.Workflow, err error) {
	workflow = &model.Workflow{}
	tx := db.GORM.Where("id = ?", id).First(&workflow)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		logger.Error("获取workflow单条数据失败," + tx.Error.Error())
		return nil, errors.New("获取workflow单条数据失败," + tx.Error.Error())
	}
	return
}

//表数据新增
func (wf *workflow) Add(workflow *model.Workflow) (err error) {
	tx := db.GORM.Create(&workflow)
	if tx.Error != nil {
		logger.Error("添加workflow数据失败," + tx.Error.Error())
		return errors.New("添加workflow数据失败," + tx.Error.Error())
	}
	return nil
}

//表数据删除
//软删除 db.GORM.Delete("id = ?", id)
//软删除执行的是UPDATE语句, 将deleted_at字段设置为时间即可, gorm默认是软删除
//实际执行语句 UPDATE `workflow` SET `deleted_at` = '2023-03-01 08:32:11' WHERE `id` IN 1
//硬删除 db.GORM.Unscoped().Delete("id = ?", id) 直接从表中删除这条数据
// 实际执行的语句 DELETE FROM 'workflow' WHERE id in 1
func (wf *workflow) DelById(id int) (err error) {
	tx := db.GORM.Where("id = ?", id).Delete(&model.Workflow{})
	if tx.Error != nil {
		logger.Error("删除workflow数据失败," + tx.Error.Error())
		return errors.New("删除workflow数据失败," + tx.Error.Error())
	}
	return nil
}