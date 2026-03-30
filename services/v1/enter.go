package v1

// ServiceGroup 服务组
type ServiceGroup struct {
	SystemUserService
	SystemRoleService
	SystemLogService
}

var ServiceGroupApp = new(ServiceGroup)
