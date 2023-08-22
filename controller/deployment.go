package controller

import (
	"fmt"
	"net/http"
	"test4/service"

	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
)

var Deployment deployment

type deployment struct{}

// 获取deployment列表, 支持过滤、排序、分页
func (d *deployment) GetDeployments(ctx *gin.Context) {
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
			"msg": "Limit/Page参数错误",
			"data": nil,
		})
		return
	}
	data, err := service.Deployment.GetDeployments(params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "获取Deployment 列表成功",
		"data": data,
	})
}

//获取deployment详情
func (d *deployment) GetDeploymentDetail(ctx *gin.Context)  {
	params := new(struct{
		DeploymentName	string	`form:"deployment_name"`
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
	data, err := service.Deployment.GetDeploymentDetail(params.DeploymentName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "获取Deployment详情成功",
		"data": data,
	})
}


//设置deployment副本数
func (d *deployment) ScaleDeployment(ctx *gin.Context)  {
	params := new(struct{
		DeploymentName  string	`json:"deployment_name"`
		Namespace		string	`json:"namespace"`
		ScaleNum		int		`json:"scale_num"`
	})
	//PUT请求, 绑定参数方法改为ctx.ShouldBindJSON
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error("Bind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.Deployment.ScaleDeployment(params.DeploymentName, params.Namespace, params.ScaleNum)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "设置Deployment副本数成功",
		"data": fmt.Sprintf("最新副本数：%d", data),
	})
}

//创建deployment
func (d *deployment) CreateDeployment(ctx *gin.Context)  {
	var (
		DeployCreate = new(service.DeployCreate)
		err error
	)
	if err = ctx.ShouldBindJSON(DeployCreate); err != nil {
		logger.Error("ShouldBind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}

	if err = service.Deployment.CreateDeployment(DeployCreate); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("创建Deployment成功:  %s", DeployCreate.Name),
		"data": nil,
	})
}

//删除deployment
func (d *deployment) DeleteDeployment(ctx *gin.Context)  {
	params := new(struct{
		DeploymentName	string	`json:"deployment_name"`
		Namespace		string	`json:"namespace"`
	})
	//DELETE请求, 绑定参数方法改为ctx.ShouldBindJSON
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error("ShouldBind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	err := service.Deployment.DeleteDeployment(params.DeploymentName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("删除Deployment %s 成功", params.DeploymentName),
		"data": nil,
	})
}

//重启deployment
func (d *deployment) RestartDeployment(ctx *gin.Context)  {
	params := new(struct{
		DeploymentName	string	`json:"deployment_name"`
		Namespace		string	`json:"namespace"`
	})
	// PUT请求, 绑定参数方法改为ctx.ShouldBindJSON
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error("ShouldBind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}

	err := service.Deployment.RestartDeployment(params.DeploymentName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("重启Deployment %s 成功 ", params.DeploymentName),
		"data": nil,
	})
}

//更新deployment
func (d *deployment) UpdateDeployment(ctx *gin.Context)  {
	params := new(struct{
		Namespace	string	`json:"namespace"`
		Content		string	`json:"content"`
	})
	fmt.Println(1111)
	// PUT请求, 绑定参数方法为ctx.SouldBindJSON
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error("ShouldBind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	err := service.Deployment.UpdateDeployment(params.Namespace, params.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "更新Deployment成功 ",
		"data": nil,
	})
}

//获取每个namespace的pod数量
func (d *deployment) GetDeployNumPerNP(ctx *gin.Context)  {
	data, err := service.Deployment.GetDeployNumPerNP()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "获取每个namespace的deployment数量成功",
		"data": data,
	})
}