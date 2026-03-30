package v1

import (
	"go-vue-admin/services/v1"
)

type ApiGroup struct {
	SystemUserApi
	SystemRoleApi
	SystemLogApi
}

var (
	systemUserService = v1.ServiceGroupApp.SystemUserService
	systemRoleService = v1.ServiceGroupApp.SystemRoleService
	systemLogService  = v1.ServiceGroupApp.SystemLogService
)

var ApiGroupApp = new(ApiGroup)
