package v1

import (
	"errors"
	"fmt"
	"strings"

	"go-vue-admin/global"
	"go-vue-admin/models"
)

type SystemRoleService struct{}

// SetRoleMenusReq 设置角色菜单权限请求
type SetRoleMenusReq struct {
	RoleID  uint   `json:"roleId" binding:"required"`
	MenuIDs []uint `json:"menuIds" binding:"required"`
}

// ==================== 角色管理 ====================

// GetRoleByID 根据ID获取角色
func (s *SystemRoleService) GetRoleByID(id uint) (*models.SystemRole, error) {
	var role models.SystemRole
	err := global.DB.First(&role, id).Error
	return &role, err
}

// GetRoleList 获取角色列表
func (s *SystemRoleService) GetRoleList(page, pageSize int, keyword string) ([]models.SystemRole, int64, error) {
	var roles []models.SystemRole
	var total int64

	db := global.DB.Model(&models.SystemRole{})

	if keyword != "" {
		db = db.Where("role_name LIKE ? OR role_code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 检查Count错误
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	err := db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&roles).Error

	return roles, total, err
}

// GetRoleOptions 获取角色选项列表（排除超级管理员，用于用户新增/编辑时的下拉选择）
func (s *SystemRoleService) GetRoleOptions() ([]models.SystemRole, error) {
	var roles []models.SystemRole
	err := global.DB.Model(&models.SystemRole{}).
		Where("role_code != ?", "admin").
		Where("status = ?", 1).
		Order("sort asc").
		Find(&roles).Error
	return roles, err
}

// CheckRoleCodeExist 检查角色代码是否已存在
func (s *SystemRoleService) CheckRoleCodeExist(roleCode string) bool {
	var count int64
	if err := global.DB.Model(&models.SystemRole{}).Where("role_code = ?", roleCode).Count(&count).Error; err != nil {
		global.Log.Errorf("检查角色代码是否存在失败: %v", err)
		return false
	}
	return count > 0
}

// CreateRole 创建角色
func (s *SystemRoleService) CreateRole(role *models.SystemRole) (uint, error) {
	if err := global.DB.Create(role).Error; err != nil {
		return 0, err
	}
	return role.ID, nil
}

// UpdateRole 更新角色
func (s *SystemRoleService) UpdateRole(role *models.SystemRole) error {
	return global.DB.Model(&models.SystemRole{}).Where("id = ?", role.ID).Updates(map[string]interface{}{
		"role_name":   role.RoleName,
		"description": role.Description,
		"status":      role.Status,
		"sort":        role.Sort,
	}).Error
}

// DeleteRole 删除角色
func (s *SystemRoleService) DeleteRole(id uint) error {
	// 检查是否有用户关联此角色（在事务外查询）
	var count int64
	if err := global.DB.Model(&models.SystemUser{}).Where("role_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return ErrRoleHasUsers
	}

	// 开启事务
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除角色的菜单权限关联
	if err := tx.Where("role_id = ?", id).Delete(&models.SystemRoleMenu{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 删除角色
	if err := tx.Delete(&models.SystemRole{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// GetRoleMenus 获取角色的菜单权限
func (s *SystemRoleService) GetRoleMenus(roleID uint) ([]uint, error) {
	var menuIDs []uint
	err := global.DB.Model(&models.SystemRoleMenu{}).Where("role_id = ?", roleID).Pluck("menu_id", &menuIDs).Error
	return menuIDs, err
}

// SetRoleMenus 设置角色的菜单权限
func (s *SystemRoleService) SetRoleMenus(req *SetRoleMenusReq) error {
	// 开启事务
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除原有的权限
	if err := tx.Where("role_id = ?", req.RoleID).Delete(&models.SystemRoleMenu{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 批量添加新的权限（优化：使用CreateInBatches）
	if len(req.MenuIDs) > 0 {
		roleMenus := make([]models.SystemRoleMenu, 0, len(req.MenuIDs))
		for _, menuID := range req.MenuIDs {
			roleMenus = append(roleMenus, models.SystemRoleMenu{
				RoleID: req.RoleID,
				MenuID: menuID,
			})
		}
		
		// 批量插入，每批100条
		if err := tx.CreateInBatches(roleMenus, 100).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	// 同步更新 Casbin 策略
	s.syncCasbinPolicy(req.RoleID, req.MenuIDs)

	return nil
}

// syncCasbinPolicy 同步 Casbin 策略
func (s *SystemRoleService) syncCasbinPolicy(roleID uint, menuIDs []uint) {
	if global.Casbin == nil {
		return
	}

	roleKey := fmt.Sprintf("role_%d", roleID)

	// 清除该角色的所有旧策略
	global.Casbin.RemoveFilteredPolicy(0, roleKey)

	// 所有角色都允许访问菜单路由接口（注意路径带 /api 前缀）
	global.Casbin.AddPolicy(roleKey, "/api/v1/system/routes", "GET")

	// 如果没有菜单权限，直接返回
	if len(menuIDs) == 0 {
		return
	}

	// 获取菜单对应的 API 路径
	var menus []models.SystemMenu
	if err := global.DB.Where("id IN ?", menuIDs).Find(&menus).Error; err != nil {
		global.Log.Errorf("获取菜单信息失败: %v", err)
		return
	}

	// 为角色添加策略（基于菜单权限，路径带 /api 前缀）
	for _, menu := range menus {
		// 根据菜单类型添加不同的权限策略
		switch menu.MenuType {
		case 1: // 目录 - 只添加查看权限
			if menu.Path != "" {
				global.Casbin.AddPolicy(roleKey, "/api/v1/system/routes", "GET")
			}
		case 2: // 菜单 - 添加查看权限
			// 根据菜单路径映射到对应的 API
			s.addMenuPolicy(roleKey, menu)
		case 3: // 按钮 - 添加操作权限
			if menu.Perm != "" {
				// 按钮权限格式: system:user:add, system:user:edit 等
				s.addButtonPolicy(roleKey, menu)
			}
		}
	}

	// 所有角色都允许访问个人中心相关接口（路径带 /api 前缀）
	global.Casbin.AddPolicy(roleKey, "/api/v1/system/users/info", "GET")
	global.Casbin.AddPolicy(roleKey, "/api/v1/system/users/profile", "PUT")
	global.Casbin.AddPolicy(roleKey, "/api/v1/system/users/password", "PUT")
}

// addMenuPolicy 添加菜单对应的 API 权限策略
func (s *SystemRoleService) addMenuPolicy(roleKey string, menu models.SystemMenu) {
	// 根据菜单路径映射到后端 API（路径带 /api 前缀）
	apiPath := s.mapMenuPathToAPI(menu.Path)
	if apiPath != "" {
		global.Casbin.AddPolicy(roleKey, "/api"+apiPath, "GET")
	}
}

// addButtonPolicy 添加按钮对应的 API 权限策略
func (s *SystemRoleService) addButtonPolicy(roleKey string, menu models.SystemMenu) {
	// 根据 perm 字段解析权限（路径带 /api 前缀）
	// perm 格式: system:user:add, system:user:edit, system:user:delete 等
	apiPath := s.mapPermToAPI(menu.Perm)
	if apiPath != "" {
		// 根据操作类型确定 HTTP 方法
		method := s.mapPermToMethod(menu.Perm)
		global.Casbin.AddPolicy(roleKey, "/api"+apiPath, method)
	}
}

// mapMenuPathToAPI 将菜单路径映射到 API 路径
func (s *SystemRoleService) mapMenuPathToAPI(menuPath string) string {
	// 简化的映射规则
	// 例如: /system/user -> /v1/system/users
	//       /system/role -> /v1/system/roles
	switch menuPath {
	case "/system/user":
		return "/v1/system/users"
	case "/system/role":
		return "/v1/system/roles"
	case "/system/menu":
		return "/v1/system/menus"
	case "/system/setting":
		return "/v1/system/settings"
	case "/system/log/operation":
		return "/v1/system/operation-logs"
	case "/system/log/login":
		return "/v1/system/login-logs"
	default:
		return ""
	}
}

// mapPermToAPI 将权限标识映射到 API 路径
func (s *SystemRoleService) mapPermToAPI(perm string) string {
	// perm 格式: system:user:add, system:user:edit 等
	// 映射到: /v1/system/users
	parts := splitPerm(perm)
	if len(parts) < 2 {
		return ""
	}

	switch parts[0] + ":" + parts[1] {
	case "system:user":
		return "/v1/system/users"
	case "system:role":
		return "/v1/system/roles"
	case "system:menu":
		return "/v1/system/menus"
	case "system:setting":
		return "/v1/system/settings"
	case "system:log:operation":
		return "/v1/system/operation-logs"
	case "system:log:login":
		return "/v1/system/login-logs"
	default:
		return ""
	}
}

// mapPermToMethod 将权限标识映射到 HTTP 方法
func (s *SystemRoleService) mapPermToMethod(perm string) string {
	// 根据 perm 后缀判断操作类型
	if strings.Contains(perm, ":add") {
		return "POST"
	}
	if strings.Contains(perm, ":edit") {
		return "PUT"
	}
	if strings.Contains(perm, ":delete") {
		return "DELETE"
	}
	if strings.Contains(perm, ":export") {
		return "GET"
	}
	return "GET"
}

// splitPerm 分割权限标识
func splitPerm(perm string) []string {
	var parts []string
	start := 0
	for i := 0; i < len(perm); i++ {
		if perm[i] == ':' {
			parts = append(parts, perm[start:i])
			start = i + 1
		}
	}
	parts = append(parts, perm[start:])
	return parts
}



// ==================== 菜单管理 ====================

// GetMenuByID 根据ID获取菜单
func (s *SystemRoleService) GetMenuByID(id uint) (*models.SystemMenu, error) {
	var menu models.SystemMenu
	err := global.DB.First(&menu, id).Error
	return &menu, err
}

// GetMenuList 获取菜单列表
func (s *SystemRoleService) GetMenuList() ([]models.SystemMenu, error) {
	var menus []models.SystemMenu
	err := global.DB.Order("sort asc").Find(&menus).Error
	return menus, err
}

// GetMenuTree 获取菜单树
func (s *SystemRoleService) GetMenuTree() ([]map[string]interface{}, error) {
	menus, err := s.GetMenuList()
	if err != nil {
		return nil, err
	}
	return s.buildMenuTree(menus, 0), nil
}

// buildMenuTree 构建菜单树
func (s *SystemRoleService) buildMenuTree(menus []models.SystemMenu, parentId uint) []map[string]interface{} {
	var tree []map[string]interface{}
	for _, menu := range menus {
		if menu.ParentID == parentId {
			// 手动构建 map，确保所有字段都正确（包括嵌套 Model 中的 ID）
			item := map[string]interface{}{
				"id":        menu.ID,
				"parentId":  menu.ParentID,
				"menuName":  menu.MenuName,
				"menuType":  menu.MenuType,
				"icon":      menu.Icon,
				"path":      menu.Path,
				"component": menu.Component,
				"perm":      menu.Perm,
				"sort":      menu.Sort,
				"status":    menu.Status,
				"visible":   menu.Visible,
				"createdAt": menu.CreatedAt,
				"updatedAt": menu.UpdatedAt,
			}
			children := s.buildMenuTree(menus, menu.ID)
			if len(children) > 0 {
				item["children"] = children
			}
			tree = append(tree, item)
		}
	}
	return tree
}

// CreateMenu 创建菜单
func (s *SystemRoleService) CreateMenu(menu *models.SystemMenu) (uint, error) {
	// 开启事务
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建菜单
	if err := tx.Create(menu).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	// 自动给超级管理员角色分配新菜单权限
	var adminRole models.SystemRole
	if err := tx.Where("role_code = ?", "admin").First(&adminRole).Error; err == nil {
		roleMenu := models.SystemRoleMenu{
			RoleID: adminRole.ID,
			MenuID: menu.ID,
		}
		if err := tx.Create(&roleMenu).Error; err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return 0, err
	}

	return menu.ID, nil
}

// UpdateMenu 更新菜单
func (s *SystemRoleService) UpdateMenu(menu *models.SystemMenu) error {
	return global.DB.Model(&models.SystemMenu{}).Where("id = ?", menu.ID).Updates(map[string]interface{}{
		"parent_id": menu.ParentID,
		"menu_name": menu.MenuName,
		"menu_type": menu.MenuType,
		"icon":      menu.Icon,
		"path":      menu.Path,
		"component": menu.Component,
		"perm":      menu.Perm,
		"sort":      menu.Sort,
		"status":    menu.Status,
		"visible":   menu.Visible,
	}).Error
}

// DeleteMenu 删除菜单
func (s *SystemRoleService) DeleteMenu(id uint) error {
	// 检查是否有子菜单
	var count int64
	if err := global.DB.Model(&models.SystemMenu{}).Where("parent_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return ErrMenuHasChildren
	}

	return global.DB.Delete(&models.SystemMenu{}, id).Error
}

// ==================== 错误定义 ====================

var (
	ErrRoleHasUsers    = errors.New("该角色下存在用户，无法删除")
	ErrMenuHasChildren = errors.New("该菜单下存在子菜单，无法删除")
)
