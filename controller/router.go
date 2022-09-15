package controller

import (
	"github.com/gin-gonic/gin"
)

//实例化router结构体，可使用对象点出首字母大写的方法（包外调用）
var Router router

//创建router结构体
type router struct {
}

//初始化路由规则，创建测试api接口

func (r *router) InitApiRouter(router *gin.Engine) {
	//router.GET("/apitest", func(ctx *gin.Context) {
	//	ctx.JSON(http.StatusOK, gin.H{
	//		"msg":  "testapi success",
	//		"data": nil,
	//	})
	//})
	router.GET("/api/k8s/pods", Pod.GetPods)
}
