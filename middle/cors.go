package middle

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//获取请求方法
		method := ctx.Request.Method

		//添加跨域请求响应头
		ctx.Header("Content-Type", "application/json")
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Max-Age", "86400")
		ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
		ctx.Header("Access-Control-Allow-Headers", "X-Token, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
		ctx.Header("Access-Control-Allow-Credentials", "false")

		//放行所有options方法
		if method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
		}

		//处理请求
		ctx.Next()
	}
}