package controller

import (
	"net/http"
	"test4/service"

	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
)

var Pod pod

type pod struct{}

/*
Controller中的方法入参是gin.Context, 用于从上下文中获取请求参数定义响应内容
流程: 绑定参数--> 调用service代码--> 根据调用结果响应具体内容
*/

// 1. 获取pod列表, 支持过滤、排序、分页
func (p *pod) GetPods(ctx *gin.Context) {
	//匿名结构体, 用于声明入参, get请求为form格式, 其他请求为json格式
	params := new(struct {
		filterName	string	`form:"filter_name"`
		Namespace 	string	`form:"namespace"`
		Page 		int		`form:"page"`
		Limit		int		`form:"limit"`
	})
	//绑定参数, 给匿名结构体中的属性赋值, 值是入参
	//form格式使用ctx.Bind方法, json格式使用ctx.ShouldBindJSON方法
	if err := ctx.Bind(params); err != nil {
		logger.Error("Bind请求参数失败, " + err.Error())
		//ctx.JSON方法用于返回响应内容, 入参是状态码和响应内容, 响应内容放入gin.H的map中
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
	//service 中的方法通过 包名.结构体.结构体变量名.方法名 使用。 service.Pod.GetPods()
	data, err := service.Pod.GetPods(params.filterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "获取Pod列表成功",
		"data": data,
	})
}

// 2.获取pod详情
func (p *pod) GetPodDetail(ctx *gin.Context)  {
	params := new(struct{
		PodName		string	`form:"pod_name"`
		Namespace	string	`form:"namespace"`
	})
	if err := ctx.Bind(params); err != nil {
		logger.Error("Bind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data":nil,
		})
		return
	}
	data, err := service.Pod.GetPodDetail(params.PodName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
			"pod_name": params.PodName,
			"namespace": params.Namespace,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "获取Pod详情成功",
		"pod_name": params.PodName,
		"data": data,
	})
}

// 3.删除pod
func (p *pod) DeletePod(ctx *gin.Context)  {
	params := new(struct{
		PodName   string `json:"pod_name"`
		Namespace string `json:"namespace"`
	})
	//Delete请求, 绑定参数方法改为ctx.ShouldBindJSON
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error("ShouldBind 请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	err := service.Pod.DeletePod(params.PodName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
			"pod_name": params.PodName,
			"namespace": params.Namespace,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "删除Pod成功",
		"data": nil,
	})
}

// 4.更新pod
func (p *pod) UpdatePod(ctx *gin.Context)  {
	params := new(struct{
		PodName		string	`json:"pod_name"`
		Namespace	string	`json:"namespace"`
		Content		string	`json:"content"`
	})
	// Put请求, 绑定参数方法改为ctx.ShouldBindJSON
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error("ShouldBind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	err := service.Pod.UpdatePod(params.PodName, params.Namespace, params.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "更新Pod成功",
		"data":nil,
	})
}

// 5. 获取pod容器
func (p *pod) GetPodContainer(ctx *gin.Context)  {
	params := new(struct{
		PodName 	string	`form:"pod_name"`
		Namespace	string	`form:"namespace"`	
	})
	//Get请求, 绑定参数方法改为ctx.Bind
	if err := ctx.Bind(params); err != nil {
		logger.Error("Bind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.Pod.GetPodContainer(params.PodName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "获取Pod容器成功",
		"data": data,
	})
}

// 6. 获取pod中容器日志
func (p *pod) GetPodLog(ctx *gin.Context) {
	params := new(struct{
		ContainerName	string	`form:"container_name"`
		PodName			string	`form:"pod_name"`
		Namespace		string	`form:"namespace"`
	})
	// Get请求, 绑定参数方法改为ctx.Bind
	if err := ctx.Bind(params); err != nil {
		logger.Error("Bind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.Pod.GetPodLog(params.ContainerName, params.PodName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "获取Pod中容器日志成功",
		"data": data,
	})
}

// 7. 获取每个namespace 的pod数量
func (p *pod) GetPodNumPerNp(ctx *gin.Context)  {
	data, err := service.Pod.GetPodNumPerNp()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "获取每个namespace的pod数量成功",
		"data": data,
	})
}