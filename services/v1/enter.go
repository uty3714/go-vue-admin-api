package v1

// ServiceGroup 服务组
type ServiceGroup struct {
	SystemUserService
	SystemRoleService
	UserService
}

var ServiceGroupApp = new(ServiceGroup)
