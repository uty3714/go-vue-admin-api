package v1

import (
	"go-vue-admin/services/v1"
)

type ApiGroup struct {
	SystemUserApi
	SystemRoleApi
	UserApi
}

var (
	systemUserService = v1.ServiceGroupApp.SystemUserService
	systemRoleService = v1.ServiceGroupApp.SystemRoleService
	userService       = v1.ServiceGroupApp.UserService
)

var ApiGroupApp = new(ApiGroup)
