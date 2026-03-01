package v1

import (
	v1 "go-vue-admin/api/v1"
	"go-vue-admin/middleware"

	"github.com/gin-gonic/gin"
)

var systemApi = v1.ApiGroupApp.SystemUserApi
var systemRoleApi = v1.ApiGroupApp.SystemRoleApi

// InitSystemRouter 初始化系统管理路由
func InitSystemRouter(rg *gin.RouterGroup) {
	router := rg.Group("/system")
	{
		// 公开路由（不需要认证）
		router.POST("/login", systemApi.Login)
		router.POST("/logout", systemApi.Logout)

		// 需要认证的路由
		authRouter := router.Use(middleware.JWTAuth())
		{
			// 动态路由菜单
			authRouter.GET("/routes", systemApi.GetAsyncRoutes)

			// 用户管理
			authRouter.GET("/user/list", systemApi.GetUserList)
			authRouter.GET("/user/info", systemApi.GetUserInfo)
			authRouter.POST("/user/create", systemApi.CreateUser)
			authRouter.PUT("/user/update", systemApi.UpdateUser)
			authRouter.DELETE("/user/delete/:id", systemApi.DeleteUser)

			// 角色管理
			authRouter.GET("/role/list", systemRoleApi.GetRoleList)
			authRouter.GET("/role/detail/:id", systemRoleApi.GetRoleDetail)
			authRouter.POST("/role/create", systemRoleApi.CreateRole)
			authRouter.PUT("/role/update", systemRoleApi.UpdateRole)
			authRouter.DELETE("/role/delete/:id", systemRoleApi.DeleteRole)
			authRouter.GET("/role/menus/:id", systemRoleApi.GetRoleMenus)
			authRouter.PUT("/role/menus", systemRoleApi.SetRoleMenus)

			// 菜单管理
			authRouter.GET("/menu/list", systemRoleApi.GetMenuList)
			authRouter.GET("/menu/tree", systemRoleApi.GetMenuTree)
			authRouter.POST("/menu/create", systemRoleApi.CreateMenu)
			authRouter.PUT("/menu/update", systemRoleApi.UpdateMenu)
			authRouter.DELETE("/menu/delete/:id", systemRoleApi.DeleteMenu)
		}
	}
}
