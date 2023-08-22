package controller

import (
	"fmt"
	"net/http"
	"test4/service"

	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
)

var DaemonSet daemonSet

type daemonSet struct{}

func (ds *daemonSet) GetDaemonSets(ctx *gin.Context) {
	params := new(struct{
		FilterName	string	`form:"filter_name"`
		Namespace	string	`form:"namespace"`
		Page		int		`form:"page"`
		Limit		int		`form:"limit"`
	})
	if err := ctx.Bind(params); err != nil {
		logger.Error("Bind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	if params.Limit <= 0 || params.Page <= 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Limit/Page参数错误或小于0",
			"data": nil,
		})
		return
	}
	data, err := service.DaemonSet.GetDaemonSets(params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "获取DaemonSet 列表成功",
		"data": data,
	})
}

func (ds *daemonSet) GetDaemonSetDetail(ctx *gin.Context)  {
	params := new(struct{
		DaemonSetName	string	`form:"daemonset_name"`
		Namespace		string	`form:"namespace"`
	})
	if err := ctx.Bind(params); err != nil {
		logger.Error("Bind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.DaemonSet.GetDaemonSetDetail(params.DaemonSetName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("获取DaemonSet: %s 详情成功", params.DaemonSetName),
		"data": data,
	})
}

func (ds *daemonSet) CreateDaemonSet(ctx *gin.Context)  {
	var (
		DaemonSetCreate = new(service.DaemonSetCreate)
		err error
	)
	if ctx.ShouldBindJSON(DaemonSetCreate); err != nil {
		logger.Error("ShouldBind请求参数绑定失败." + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	err = service.DaemonSet.CreateDaemonSet(DaemonSetCreate)
	if err != nil {
		logger.Error("创建DaemonSet失败." + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "DaemonSet创建成功.",
		"data": nil,
	})
}

func (ds *daemonSet) DeleteDaemonSet(ctx *gin.Context)  {
	params := new(struct{
		DaemonSetName	string	`json:"daemonset_name"`
		Namespace		string	`json:"namespace"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error("ShouldBind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	err := service.DaemonSet.DeleteDaemonSet(params.DaemonSetName, params.Namespace)
	if err != nil {
		logger.Error("删除DaemonSet失败." + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("删除DaemonSet: %s 成功", params.DaemonSetName), 
		"data": nil,
	})
}

func (ds *daemonSet) RestartDaemonSet(ctx *gin.Context)  {
	params := new(struct{
		DaemonSetName	string	`json:"daemonset_name"`
		Namespace		string	`json:"namespace"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error("ShouldBind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	err := service.DaemonSet.RestartDaemonSet(params.DaemonSetName, params.Namespace)
	if err != nil {
		logger.Error("重启DaemonSet失败." + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("重启DaemonSet: %s 成功", params.DaemonSetName), 
		"data": nil,
	})
}

func (ds *daemonSet) UpdateDaemonSet(ctx *gin.Context)  {
	params := new(struct{
		Content			string	`json:"content"`
		Namespace		string	`json:"namespace"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error("ShouldBind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	err := service.DaemonSet.UpdateDaemonSet(params.Namespace, params.Content)
	if err != nil {
		logger.Error("更新DaemonSet失败." + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "更新DaemonSet成功", 
		"data": nil,
	})
}

func (ds *daemonSet) GetDaemonSetNumPerNp(ctx *gin.Context)  {
	data, err := service.DaemonSet.GetDaemonSetNumPerNp()
	if err != nil {
		logger.Error("获取各个namespace下的DaemonSet数量失败." + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "获取各个namespace下的DaemonSet数量成功.", 
		"data": data,
	})
}