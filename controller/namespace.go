package controller

import (
	"errors"
	"fmt"
	"net/http"
	"test4/service"

	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
)

var Namepsace namespace

type namespace struct{}

func (ns *namespace) GetNamespaces(ctx *gin.Context) {
	params := new(struct{
		FilterName	string	`form:"filter_name"`
		Limit		int		`form:"limit"`
		Page		int		`form:"page"`
	})
	if err := ctx.Bind(params); err != nil {
		logger.Error(errors.New("Bind请求参数绑定失败. " + err.Error()))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	if params.Limit <=0 || params.Page <=0 {
		logger.Error(errors.New("Limit/Page 参数不合法或小于等于0 "))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Limit/Page 参数不合法或小于等于0 ",
			"data": nil,
		})
		return
	}
	data, err := service.Namepsace.GetNamespaces(params.FilterName, params.Limit, params.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "获取namespace列表成功",
		"data": data,
	})
}

func (ns *namespace) GetNamespaceDetail(ctx *gin.Context)  {
	params := new(struct{
		NamespaceName	string	`form:"namespace_name"`
	})
	if err := ctx.Bind(params); err != nil {
		logger.Error(errors.New("Bind请求参数绑定失败. " + err.Error()))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.Namepsace.GetNamespaceDetail(params.NamespaceName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("获取Namespace: %s 详情成功", params.NamespaceName),
		"data": data,
	})
}

func (ns *namespace) DeleteNamespace(ctx *gin.Context)  {
	params := new(struct{
		NamespaceName	string	`json:"namespace_name"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error(errors.New("ShouldBind请求参数绑定失败. " + err.Error()))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	err := service.Namepsace.DeleteNamespace(params.NamespaceName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("删除Namespace: %s 成功", params.NamespaceName),
		"data": nil,
	})
}