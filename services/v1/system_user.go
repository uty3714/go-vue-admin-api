package v1

import (
	"errors"
	"go-vue-admin/global"
	"go-vue-admin/models"
	"go-vue-admin/models/res"
	"go-vue-admin/util"
	"time"
)

type SystemUserService struct{}

// ==================== 基础 CRUD ====================

// GetUserByID 根据ID获取用户
func (s *SystemUserService) GetUserByID(id uint) (*models.SystemUser, error) {
	var user models.SystemUser
	err := global.DB.Preload("Role").First(&user, id).Error
	return &user, err
}

// GetUserByUsername 根据用户名获取用户
func (s *SystemUserService) GetUserByUsername(username string) (*models.SystemUser, error) {
	var user models.SystemUser
	err := global.DB.Preload("Role").Where("username = ?", username).First(&user).Error
	return &user, err
}

// CheckUserExist 检查用户名是否已存在
func (s *SystemUserService) CheckUserExist(username string) bool {
	var count int64
	global.DB.Model(&models.SystemUser{}).Where("username = ?", username).Count(&count)
	return count > 0
}

// CheckUserExistExceptID 检查用户名是否已存在（排除指定ID）
func (s *SystemUserService) CheckUserExistExceptID(username string, excludeID uint) bool {
	var count int64
	global.DB.Model(&models.SystemUser{}).Where("username = ? AND id != ?", username, excludeID).Count(&count)
	return count > 0
}

// CreateUser 创建用户
func (s *SystemUserService) CreateUser(req *models.SystemUserReq) (uint, error) {
	user := models.SystemUser{
		Username: req.Username,
		Password: util.BcryptHash(req.Password),
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		Status:   req.Status,
		RoleID:   req.RoleID,
	}
	if err := global.DB.Create(&user).Error; err != nil {
		return 0, err
	}
	return user.ID, nil
}

// UpdateUser 更新用户
func (s *SystemUserService) UpdateUser(req *models.SystemUserUpdateReq) error {
	var user models.SystemUser
	if err := global.DB.First(&user, req.ID).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{}

	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Status != 0 {
		updates["status"] = req.Status
	}
	if req.RoleID != 0 {
		updates["role_id"] = req.RoleID
	}
	if req.Password != "" {
		updates["password"] = util.BcryptHash(req.Password)
	}

	if len(updates) == 0 {
		return nil
	}

	return global.DB.Model(&user).Updates(updates).Error
}

// DeleteUser 删除用户
func (s *SystemUserService) DeleteUser(id uint) error {
	return global.DB.Delete(&models.SystemUser{}, id).Error
}

// GetUserList 获取用户列表
func (s *SystemUserService) GetUserList(page, pageSize int, keyword, status string) ([]models.SystemUser, int64, error) {
	var users []models.SystemUser
	var total int64

	db := global.DB.Model(&models.SystemUser{}).Preload("Role")

	if keyword != "" {
		db = db.Where("username LIKE ? OR nickname LIKE ? OR phone LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status != "" {
		db = db.Where("status = ?", util.StringToInt(status))
	}

	db.Count(&total)
	err := db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error

	// 隐藏密码
	for i := range users {
		users[i].Password = ""
	}

	return users, total, err
}

// ==================== 业务方法 ====================

// Login 用户登录
func (s *SystemUserService) Login(req *models.SystemUserLoginReq, clientIP string) (*models.SystemUserLoginRes, int) {
	// 查询用户（预加载角色信息）
	var user models.SystemUser
	if err := global.DB.Preload("Role").Where("username = ?", req.Username).First(&user).Error; err != nil {
		return nil, res.ErrorCodeLoginFailed
	}

	// 验证密码
	if !util.BcryptCheck(req.Password, user.Password) {
		return nil, res.ErrorCodePasswordError
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, res.ErrorCodeUserDisabled
	}

	// 生成token
	j := util.NewJWT()
	claims := j.CreateClaims(util.CustomClaims{
		UserID:   user.ID,
		Username: user.Username,
		RoleID:   user.RoleID,
	})
	token, err := j.CreateToken(claims)
	if err != nil {
		return nil, res.ErrorCodeInternalServer
	}

	// 更新登录信息
	now := time.Now().Format("2006-01-02 15:04:05")
	global.DB.Model(&user).Updates(map[string]interface{}{
		"last_login_ip": clientIP,
		"last_login_at": now,
	})

	// 隐藏密码
	user.Password = ""
	// 设置前端需要的 roles 字段
	user.Roles = []string{user.Role.RoleCode}

	return &models.SystemUserLoginRes{
		Token:     token,
		ExpiresAt: claims.ExpiresAt.Unix(),
		UserInfo:  user,
	}, res.SuccessCode
}

// GetUserInfo 获取当前用户信息
func (s *SystemUserService) GetUserInfo(userID uint) (*models.SystemUser, error) {
	user, err := s.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	user.Password = ""
	user.Roles = []string{user.Role.RoleCode}
	return user, nil
}

// GetAsyncRoutes 获取当前用户的动态路由菜单
func (s *SystemUserService) GetAsyncRoutes(userID uint) ([]map[string]interface{}, int) {
	// 查询用户信息及角色
	user, err := s.GetUserByID(userID)
	if err != nil {
		return nil, res.ErrorCodeUserNotExist
	}

	var menus []models.SystemMenu

	// 查询角色的菜单权限（所有角色都根据权限表查询，包括超级管理员）
	var menuIDs []uint
	global.DB.Model(&models.SystemRoleMenu{}).Where("role_id = ?", user.RoleID).Pluck("menu_id", &menuIDs)

	// 查询菜单详情
	if len(menuIDs) > 0 {
		global.DB.Where("id IN ?", menuIDs).Where("status = ?", 1).Where("visible = ?", 1).Order("sort asc").Find(&menus)
	}

	// 转换为前端路由格式
	routes := s.buildRoutesFromMenus(menus, user.Role.RoleCode)

	return routes, res.SuccessCode
}

// buildRoutesFromMenus 将菜单列表转换为前端路由格式
func (s *SystemUserService) buildRoutesFromMenus(menus []models.SystemMenu, roleCode string) []map[string]interface{} {
	var routes []map[string]interface{}

	// 构建菜单树
	menuTree := s.buildMenuTreeForRoutes(menus, 0)

	// 转换为前端路由格式
	for _, menu := range menuTree {
		route := s.menuToRoute(menu, roleCode)
		if route != nil {
			routes = append(routes, route)
		}
	}

	return routes
}

// buildMenuTreeForRoutes 构建菜单树（用于路由）
func (s *SystemUserService) buildMenuTreeForRoutes(menus []models.SystemMenu, parentId uint) []map[string]interface{} {
	var tree []map[string]interface{}
	for _, menu := range menus {
		if menu.ParentID == parentId {
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
			}
			children := s.buildMenuTreeForRoutes(menus, menu.ID)
			if len(children) > 0 {
				item["children"] = children
			}
			tree = append(tree, item)
		}
	}
	return tree
}

// menuToRoute 将菜单转换为前端路由格式
func (s *SystemUserService) menuToRoute(menu map[string]interface{}, roleCode string) map[string]interface{} {
	menuType, _ := menu["menuType"].(int)

	// 按钮不生成路由
	if menuType == 3 {
		return nil
	}

	path, _ := menu["path"].(string)
	menuName, _ := menu["menuName"].(string)
	icon, _ := menu["icon"].(string)
	component, _ := menu["component"].(string)

	// 构建meta
	meta := map[string]interface{}{
		"title": menuName,
		"icon":  icon,
	}

	// 目录和菜单都显示在侧边栏
	meta["showLink"] = true

	// 添加角色权限（当前用户的角色）
	meta["roles"] = []string{roleCode}

	// 构建路由
	route := map[string]interface{}{
		"path": path,
		"meta": meta,
	}

	// 设置name（目录加Parent后缀，避免和子菜单冲突）
	if menuType == 1 {
		route["name"] = menuName + "Parent"
	} else {
		route["name"] = menuName
	}

	// 设置组件
	if component != "" {
		route["component"] = component
	}

	// 处理子路由
	if children, ok := menu["children"].([]map[string]interface{}); ok && len(children) > 0 {
		var childRoutes []map[string]interface{}
		for _, child := range children {
			childRoute := s.menuToRoute(child, roleCode)
			if childRoute != nil {
				childRoutes = append(childRoutes, childRoute)
			}
		}
		if len(childRoutes) > 0 {
			route["children"] = childRoutes
			// 目录添加redirect指向第一个子菜单
			if menuType == 1 {
				firstChildPath, _ := childRoutes[0]["path"].(string)
				route["redirect"] = firstChildPath
			}
		}
	}

	return route
}

// ==================== 内部辅助方法 ====================

// 为了兼容现有错误处理
var ErrUserNotExist = errors.New("用户不存在")
var ErrInvalidPassword = errors.New("密码错误")
var ErrUserDisabled = errors.New("用户已被禁用")
