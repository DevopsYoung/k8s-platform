package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"k8s-platform/service"
	"net/http"
)

var Pod pod

type pod struct {
}

//获取pod列表，支持分页、过滤、排序
func (p *pod) GetPods(ctx *gin.Context) {
	//处理入参
	//匿名结构体用来定义入参，get请求为form格式，其他请求为json格式
	params := new(struct {
		FilterName string `form:"filter_name"`
		NameSpace  string `form:"namespace"`
		Limit      int    `form:"limit"`
		Page       int    `form:"page"`
	})

	//form格式使用Bind方法，json格式使用ShouldBindJson方法
	if err := ctx.Bind(params); err != nil {
		logger.Error("参数绑定失败" + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "参数绑定失败" + err.Error(),
			"data": nil,
		})
	}
	data, err := service.Pod.GetPods(params.FilterName, params.NameSpace, params.Limit, params.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"mst":  err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取pod列表成功",
		"data": data,
	})
}
