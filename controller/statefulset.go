package controller

import (
	"fmt"
	"net/http"
	"test4/service"

	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
)

var StatefulSet statefulSet

type statefulSet struct{}

func (s *statefulSet) GetStatefulSets(ctx *gin.Context) {
	params := new(struct{
		FilterName	string	`form:"filter_name"`
		Namespace	string	`form:"namespace"`
		Limit		int		`form:"limit"`
		Page		int		`form:"page"`
	})
	if err := ctx.Bind(params); err != nil {
		logger.Error("Bind请求参数绑定失败. " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	if params.Limit <=0 || params.Page <=0 {
		logger.Error("Limit/Page 参数不合法或小于等于0 ")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Limit/Page 参数不合法或小于等于0 ",
			"data": nil,
		})
		return
	}
	data, err := service.StatefulSet.GetStatefulSets(params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "statefulset列表获取成功",
		"data": data,
	})
}

func (s *statefulSet) GetStatefulSetDetail(ctx *gin.Context)  {
	params := new(struct{
		StatefulSetName		string	`form:"statefulset_name"`
		Namespace			string	`form:"namespace"`
	})
	if err := ctx.Bind(params); err != nil {
		logger.Error("Bind请求参数绑定失败. " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.StatefulSet.GetStatefulSetDetail(params.StatefulSetName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("statefulset: %s列表获取成功", params.StatefulSetName),
		"data": data,
	})
}

func (s *statefulSet) ScaleStatefulSet(ctx *gin.Context)  {
	params := new(struct{
		StatefulSetName		string	`json:"statefulset_name"`
		Namespace			string	`json:"namespace"`
		ScaleNum			int		`json:"scale_num"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error("ShouldBind请求参数绑定失败. " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.StatefulSet.ScaleStatefulSet(params.StatefulSetName, params.Namespace, params.ScaleNum)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("statefulset: %s副本数更新成功", params.StatefulSetName),
		"data": fmt.Sprintf("statefulset: %s副本数更新成功, 最新副本数为：%d", params.StatefulSetName, data),
	})
}

func (s *statefulSet) CreateStatefulSet(ctx *gin.Context)  {
	var (
		StatefulSetCreate = new(service.StatefulSetCreate)
		err error
	)
	if err := ctx.ShouldBindJSON(StatefulSetCreate); err != nil {
		logger.Error("ShouldBind请求参数绑定失败. " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}

	err = service.StatefulSet.CreateStatefulSet(StatefulSetCreate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("statefulset: %s 创建成功", StatefulSetCreate.StatefulSetName),
		"data": fmt.Sprintf("statefulset: %s 创建成功", StatefulSetCreate.StatefulSetName),
	})
}

func (s *statefulSet) DeleteStatefulSet(ctx *gin.Context)  {
	params := new(struct{
		StatefulSetName		string	`json:"statefulset_name"`
		Namespace			string	`json:"namespace"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error("ShouldBind请求参数绑定失败. " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	err := service.StatefulSet.DeleteStatefulSet(params.StatefulSetName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("statefulset: %s 删除成功", params.StatefulSetName),
		"data": fmt.Sprintf("statefulset: %s 删除成功", params.StatefulSetName),
	})
}


func (s *statefulSet) RestartStatefulSet(ctx *gin.Context)  {
	params := new(struct{
		StatefulSetName		string	`json:"statefulset_name"`
		Namespace			string	`json:"namespace"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error("ShouldBind请求参数绑定失败. " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	err := service.StatefulSet.RestartStatefulSet(params.StatefulSetName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("statefulset: %s 重启成功", params.StatefulSetName),
		"data": fmt.Sprintf("statefulset: %s 重启成功", params.StatefulSetName),
	})
}

func (s *statefulSet) UpdateStatefulSet(ctx *gin.Context)  {
	params := new(struct{
		Content		string	`json:"content"`
		Namespace	string	`json:"namespace"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error("ShouldBind请求参数绑定失败. " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	err := service.StatefulSet.UpdateStatefulSet(params.Namespace, params.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "statefulset更新成功",
		"data": "statefulset更新成功",
	})
}

func (s *statefulSet) GetStatefulSetsNumPerNp(ctx *gin.Context) {
	data, err := service.StatefulSet.GetStatefulSetsNumPerNp()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "获取每个namespace下的statefulset成功.",
		"data": data,
	})
}