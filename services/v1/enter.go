package v1

// ServiceGroup 服务组
type ServiceGroup struct {
	SystemUserService
	SystemRoleService
	SystemLogService
	SystemSettingService
}

var ServiceGroupApp = new(ServiceGroup)
