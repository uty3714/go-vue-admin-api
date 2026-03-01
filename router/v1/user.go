package v1

import (
	v1 "go-vue-admin/api/v1"
	"go-vue-admin/middleware"

	"github.com/gin-gonic/gin"
)

var userApi = v1.ApiGroupApp.UserApi

// InitUserRouter 初始化用户路由
func InitUserRouter(rg *gin.RouterGroup) {
	router := rg.Group("/user")
	{
		// 小程序端路由
		miniRouter := router.Use(middleware.MiniAuth())
		{
			miniRouter.GET("/info", userApi.GetUserInfo)
			miniRouter.PUT("/update", userApi.UpdateUser)
			miniRouter.GET("/address/list", userApi.GetAddressList)
			miniRouter.POST("/address/create", userApi.CreateAddress)
			miniRouter.PUT("/address/update/:id", userApi.UpdateAddress)
			miniRouter.DELETE("/address/delete/:id", userApi.DeleteAddress)
			miniRouter.GET("/cart/list", userApi.GetCartList)
			miniRouter.POST("/cart/add", userApi.AddCart)
			miniRouter.PUT("/cart/update/:id", userApi.UpdateCart)
			miniRouter.DELETE("/cart/delete/:id", userApi.DeleteCart)
		}
	}
}
