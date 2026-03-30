package models

// SystemUser 系统用户表
type SystemUser struct {
	ID          uint       `gorm:"primarykey;comment:主键ID" json:"id"`
	CreatedAt   LocalTime  `gorm:"type:datetime;not null;comment:创建时间" json:"createdAt"`
	UpdatedAt   LocalTime  `gorm:"type:datetime;not null;comment:更新时间" json:"updatedAt"`
	Username    string     `gorm:"type:varchar(64);not null;uniqueIndex;comment:用户名" json:"username"`
	Password    string     `gorm:"type:varchar(128);not null;comment:密码" json:"-"`
	Nickname    string     `gorm:"type:varchar(64);comment:昵称" json:"nickname"`
	Avatar      string     `gorm:"type:varchar(255);comment:头像URL" json:"avatar"`
	Email       string     `gorm:"type:varchar(128);comment:邮箱" json:"email"`
	Phone       string     `gorm:"type:varchar(20);comment:手机号" json:"phone"`
	Status      int        `gorm:"type:tinyint;default:1;comment:状态 1启用 2禁用" json:"status"`
	RoleID      uint       `gorm:"index;comment:角色ID" json:"roleId"`
	Role        SystemRole `gorm:"foreignKey:RoleID" json:"role"`
	LastLoginIP string     `gorm:"type:varchar(128);default:null;comment:最后登录IP" json:"lastLoginIp"`
	LastLoginAt *string    `gorm:"type:datetime;default:null;comment:最后登录时间" json:"lastLoginAt"`
	Roles       []string   `gorm:"-" json:"roles"` // 前端需要的角色数组格式
}

func (SystemUser) TableName() string {
	return "system_user"
}

// SystemUserReq 系统用户请求参数（创建时使用）
type SystemUserReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	RoleID   uint   `json:"roleId"`
	Status   int    `json:"status"`
}

// SystemUserUpdateReq 系统用户更新请求参数（更新时使用，字段都是可选的）
type SystemUserUpdateReq struct {
	ID       uint   `json:"id" binding:"required"`
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	RoleID   uint   `json:"roleId"`
	Status   int    `json:"status"`
}

// SystemUserLoginReq 登录请求参数
type SystemUserLoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// SystemUserLoginRes 登录响应参数
type SystemUserLoginRes struct {
	Token     string      `json:"token"`
	ExpiresAt int64       `json:"expiresAt"`
	UserInfo  SystemUser  `json:"userInfo"`
}

// SystemUserProfileReq 当前用户更新个人信息请求参数
 type SystemUserProfileReq struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

// SystemUserPasswordReq 当前用户修改密码请求参数
type SystemUserPasswordReq struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}
