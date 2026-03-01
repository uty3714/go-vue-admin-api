package models

// SystemRole 系统角色表
type SystemRole struct {
	Model
	RoleName    string `gorm:"type:varchar(128);not null;comment:角色名称" json:"roleName"`
	RoleCode    string `gorm:"type:varchar(128);not null;uniqueIndex;comment:角色代码" json:"roleCode"`
	Description string `gorm:"type:varchar(255);comment:角色描述" json:"description"`
	Status      int    `gorm:"type:tinyint;default:1;comment:状态 1启用 2禁用" json:"status"`
	Sort        int    `gorm:"type:int;default:0;comment:排序" json:"sort"`
}

func (SystemRole) TableName() string {
	return "system_role"
}

// SystemRoleMenu 角色菜单关联表
type SystemRoleMenu struct {
	ID     uint `gorm:"primarykey" json:"id"`
	RoleID uint `gorm:"not null;index;comment:角色ID" json:"roleId"`
	MenuID uint `gorm:"not null;index;comment:菜单ID" json:"menuId"`
}

func (SystemRoleMenu) TableName() string {
	return "system_role_menu"
}

// SystemMenu 系统菜单表
type SystemMenu struct {
	Model
	ParentID  uint   `gorm:"index;default:0;comment:父菜单ID" json:"parentId"`
	MenuName  string `gorm:"type:varchar(128);not null;comment:菜单名称" json:"menuName"`
	MenuType  int    `gorm:"type:tinyint;default:1;comment:菜单类型 1目录 2菜单 3按钮" json:"menuType"`
	Icon      string `gorm:"type:varchar(128);comment:菜单图标" json:"icon"`
	Path      string `gorm:"type:varchar(255);comment:路由路径" json:"path"`
	Component string `gorm:"type:varchar(255);comment:组件路径" json:"component"`
	Perm      string `gorm:"type:varchar(255);comment:权限标识" json:"perm"`
	Sort      int    `gorm:"type:int;default:0;comment:排序" json:"sort"`
	Status    int    `gorm:"type:tinyint;default:1;comment:状态 1启用 2禁用" json:"status"`
	Visible   int    `gorm:"type:tinyint;default:1;comment:是否显示 1是 2否" json:"visible"`
}

func (SystemMenu) TableName() string {
	return "system_menu"
}
