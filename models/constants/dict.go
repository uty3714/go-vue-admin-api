package constants

// ==================== 用户状态 ====================
const (
	UserStatusEnabled  = 1 // 启用
	UserStatusDisabled = 2 // 禁用
)

var UserStatusMap = map[int]string{
	UserStatusEnabled:  "启用",
	UserStatusDisabled: "禁用",
}

// ==================== 角色状态 ====================
const (
	RoleStatusEnabled  = 1 // 启用
	RoleStatusDisabled = 2 // 禁用
)

// ==================== 菜单类型 ====================
const (
	MenuTypeDirectory = 1 // 目录
	MenuTypeMenu      = 2 // 菜单
	MenuTypeButton    = 3 // 按钮
)

var MenuTypeMap = map[int]string{
	MenuTypeDirectory: "目录",
	MenuTypeMenu:      "菜单",
	MenuTypeButton:    "按钮",
}

// ==================== 菜单状态 ====================
const (
	MenuStatusEnabled  = 1 // 启用
	MenuStatusDisabled = 2 // 禁用
)

// ==================== 菜单可见性 ====================
const (
	MenuVisibleShow = 1 // 显示
	MenuVisibleHide = 2 // 隐藏
)

// ==================== 操作日志状态 ====================
const (
	OperationLogStatusSuccess = 1 // 成功
	OperationLogStatusFailed  = 2 // 失败
)

// ==================== 登录日志状态 ====================
const (
	LoginLogStatusSuccess = 1 // 成功
	LoginLogStatusFailed  = 2 // 失败
)

// ==================== HTTP方法 ====================
const (
	HTTPMethodGet     = "GET"
	HTTPMethodPost    = "POST"
	HTTPMethodPut     = "PUT"
	HTTPMethodDelete  = "DELETE"
	HTTPMethodPatch   = "PATCH"
	HTTPMethodOptions = "OPTIONS"
)

// ==================== 系统角色代码 ====================
const (
	RoleCodeAdmin = "admin" // 超级管理员
	RoleCodeUser  = "user"  // 普通用户
)

// ==================== JWT相关 ====================
const (
	JWTTokenTypeBearer = "Bearer"
)

// ==================== 响应码 ====================
const (
	ResponseCodeSuccess = 200
	ResponseCodeError   = 500
)
