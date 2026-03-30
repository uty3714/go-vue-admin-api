package v1

import (
	"go-vue-admin/middleware"

	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由
func InitRouter(r *gin.Engine) *gin.RouterGroup {
	// 全局中间件
	r.Use(middleware.Cors())

	// API v1 版本路由组
	apiV1 := r.Group("/api/v1")
	{
		// 系统管理路由
		InitSystemRouter(apiV1)
	}

	return apiV1
}
