package flag

import (
	"fmt"
	"go-vue-admin/global"
	"go-vue-admin/models"
	"go-vue-admin/util"
	"os"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

// ResetAdminPassword 重置管理员密码为 admin123
func ResetAdminPassword() {
	db := global.DB
	if db == nil {
		fmt.Println("数据库连接失败")
		return
	}

	// 查找管理员账号
	var user models.SystemUser
	if err := db.Where("username = ?", "admin").First(&user).Error; err != nil {
		fmt.Println("未找到管理员账号(admin)")
		return
	}

	// 生成新密码哈希
	newPassword := "admin123"
	hash := util.BcryptHash(newPassword)

	// 更新密码
	if err := db.Model(&user).Update("password", hash).Error; err != nil {
		fmt.Printf("密码重置失败: %v\n", err)
		return
	}

	fmt.Println("✅ 管理员密码已重置为: admin123")
}

// ResetDB 重置数据库（删除所有表并重新创建）
func ResetDB() {
	fmt.Println("⚠️  警告: 即将删除所有数据表并重新初始化！")
	fmt.Println("开始重置数据库...")

	db := global.DB
	if db == nil {
		fmt.Println("数据库连接失败")
		return
	}

	// 删除所有表（按依赖关系倒序删除）
	tables := []interface{}{
		&models.OperationLog{},
		&models.LoginLog{},
		&models.SystemRoleMenu{},
		&models.SystemMenu{},
		&models.SystemUser{},
		&models.SystemRole{},
	}

	for _, table := range tables {
		if err := db.Migrator().DropTable(table); err != nil {
			fmt.Printf("删除表失败: %v\n", err)
		}
	}

	// 删除 Casbin 规则表
	db.Migrator().DropTable("casbin_rule")

	fmt.Println("所有数据表已删除")

	// 重新初始化
	MigrateDB()
}

// MigrateDB 数据库迁移
func MigrateDB() {
	fmt.Println("开始初始化数据库...")

	db := global.DB
	if db == nil {
		fmt.Println("数据库连接失败")
		return
	}

	// 设置数据库字符集为 utf8mb4
	db.Exec("SET NAMES utf8mb4")
	db.Exec("SET CHARACTER SET utf8mb4")
	db.Exec("SET character_set_connection=utf8mb4")

	// 自动迁移表结构
	err := db.AutoMigrate(
		&models.SystemUser{},
		&models.SystemRole{},
		&models.SystemRoleMenu{},
		&models.SystemMenu{},
		&models.SystemSetting{},
		&models.OperationLog{},
		&models.LoginLog{},
	)
	if err != nil {
		fmt.Printf("数据库迁移失败: %v\n", err)
		return
	}

	// 修复表字符集
	fixCharset()

	fmt.Println("数据库迁移完成")

	// 初始化基础数据
	initBaseData()

	fmt.Println("数据库初始化完成")
}

// fixCharset 修复所有表的字符集为 utf8mb4
func fixCharset() {
	fmt.Println("开始修复表字符集...")
	db := global.DB
	tables := []string{
		"system_user",
		"system_role",
		"system_role_menu",
		"system_menu",
		"system_setting",
		"system_operation_log",
		"system_login_log",
	}
	
	for _, table := range tables {
		sql := fmt.Sprintf("ALTER TABLE %s CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", table)
		if err := db.Exec(sql).Error; err != nil {
			fmt.Printf("  修复表 %s 字符集失败: %v\n", table, err)
		} else {
			fmt.Printf("  表 %s 字符集已设置为 utf8mb4\n", table)
		}
	}
	
	// 修复长文本列字符集
	columns := []struct {
		table  string
		column string
		ctype  string
	}{
		{"system_operation_log", "response_data", "LONGTEXT"},
		{"system_operation_log", "request_data", "LONGTEXT"},
		{"system_operation_log", "error_message", "TEXT"},
		{"system_login_log", "message", "VARCHAR(255)"},
	}
	
	for _, col := range columns {
		sql := fmt.Sprintf("ALTER TABLE %s MODIFY %s %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", 
			col.table, col.column, col.ctype)
		if err := db.Exec(sql).Error; err != nil {
			fmt.Printf("  修复列 %s.%s 字符集失败: %v\n", col.table, col.column, err)
		} else {
			fmt.Printf("  列 %s.%s 字符集已设置为 utf8mb4\n", col.table, col.column)
		}
	}
	
	fmt.Println("表字符集修复完成")
}

// initBaseData 初始化基础数据
func initBaseData() {
	db := global.DB

	// 检查是否已存在角色数据
	var count int64
	db.Model(&models.SystemRole{}).Count(&count)
	if count > 0 {
		fmt.Println("基础数据已存在，跳过初始化")
		// 但菜单数据可能需要初始化（兼容旧数据）
		initMenuData()
		// 同步所有角色的Casbin策略
		fmt.Println("正在同步所有角色的权限策略...")
		initCasbinPolicyForAllRoles()
		return
	}

	// 创建超级管理员角色
	adminRole := models.SystemRole{
		RoleName:    "超级管理员",
		RoleCode:    "admin",
		Description: "系统超级管理员，拥有所有权限",
		Status:      1,
		Sort:        0,
	}
	if err := db.Create(&adminRole).Error; err != nil {
		fmt.Printf("创建管理员角色失败: %v\n", err)
		return
	}

	// 创建普通用户角色
	userRole := models.SystemRole{
		RoleName:    "普通用户",
		RoleCode:    "user",
		Description: "普通用户角色",
		Status:      1,
		Sort:        1,
	}
	if err := db.Create(&userRole).Error; err != nil {
		fmt.Printf("创建用户角色失败: %v\n", err)
		return
	}

	// 创建默认管理员账号
	adminUser := models.SystemUser{
		Username:    "admin",
		Password:    util.BcryptHash("admin123"),
		Nickname:    "管理员",
		Email:       "admin@go-vue-admin.com",
		Phone:       "13800138000",
		Status:      1,
		RoleID:      adminRole.ID,
		LastLoginIP: "",
		LastLoginAt: nil,
	}
	if err := db.Create(&adminUser).Error; err != nil {
		fmt.Printf("创建管理员账号失败: %v\n", err)
		return
	}

	fmt.Printf("基础数据初始化完成\n")
	fmt.Printf("管理员账号: admin, 密码: admin123\n")

	// 初始化菜单数据及权限
	initMenuData()

	// 初始化系统设置（默认关闭日志）
	initSystemSetting()

	// 初始化Casbin权限
	initCasbinPolicy(adminRole.ID)
}

// initMenuData 初始化菜单数据
func initMenuData() {
	db := global.DB

	// 检查是否已存在菜单数据
	var count int64
	db.Model(&models.SystemMenu{}).Count(&count)
	if count > 0 {
		fmt.Println("菜单数据已存在，跳过初始化")
		return
	}

	// 注意：首页(welcome)使用前端静态路由，不在后端返回

	// 创建系统管理目录
	systemDir := models.SystemMenu{
		ParentID:  0,
		MenuName:  "系统管理",
		MenuType:  1, // 目录
		Icon:      "ri:settings-3-line",
		Path:      "/system",
		Component: "",
		Perm:      "system:view",
		Sort:      1,
		Status:    1,
		Visible:   1,
	}
	if err := db.Create(&systemDir).Error; err != nil {
		fmt.Printf("创建系统管理目录失败: %v\n", err)
		return
	}

	// 创建用户管理菜单
	userMenu := models.SystemMenu{
		ParentID:  systemDir.ID,
		MenuName:  "用户管理",
		MenuType:  2, // 菜单
		Icon:      "ri:admin-line",
		Path:      "/system/user",
		Component: "system/user/index",
		Perm:      "system:user:view",
		Sort:      1,
		Status:    1,
		Visible:   1,
	}
	if err := db.Create(&userMenu).Error; err != nil {
		fmt.Printf("创建用户管理菜单失败: %v\n", err)
		return
	}

	// 创建角色管理菜单
	roleMenu := models.SystemMenu{
		ParentID:  systemDir.ID,
		MenuName:  "角色管理",
		MenuType:  2, // 菜单
		Icon:      "ri:shield-keyhole-line",
		Path:      "/system/role",
		Component: "system/role/index",
		Perm:      "system:role:view",
		Sort:      2,
		Status:    1,
		Visible:   1,
	}
	if err := db.Create(&roleMenu).Error; err != nil {
		fmt.Printf("创建角色管理菜单失败: %v\n", err)
		return
	}

	// 创建菜单管理菜单
	menuMenu := models.SystemMenu{
		ParentID:  systemDir.ID,
		MenuName:  "菜单管理",
		MenuType:  2, // 菜单
		Icon:      "ep:menu",
		Path:      "/system/menu",
		Component: "system/menu/index",
		Perm:      "system:menu:view",
		Sort:      3,
		Status:    1,
		Visible:   1,
	}
	if err := db.Create(&menuMenu).Error; err != nil {
		fmt.Printf("创建菜单管理菜单失败: %v\n", err)
		return
	}

	// 创建系统设置菜单
	settingMenu := models.SystemMenu{
		ParentID:  systemDir.ID,
		MenuName:  "系统设置",
		MenuType:  2, // 菜单
		Icon:      "ri:settings-4-line",
		Path:      "/system/setting",
		Component: "system/setting/index",
		Perm:      "system:setting:view",
		Sort:      4,
		Status:    1,
		Visible:   1,
	}
	if err := db.Create(&settingMenu).Error; err != nil {
		fmt.Printf("创建系统设置菜单失败: %v\n", err)
		return
	}

	// 创建操作日志菜单
	operationLogMenu := models.SystemMenu{
		ParentID:  systemDir.ID,
		MenuName:  "操作日志",
		MenuType:  2, // 菜单
		Icon:      "ri:file-list-line",
		Path:      "/system/log/operation",
		Component: "system/log/operation",
		Perm:      "system:log:operation:view",
		Sort:      5,
		Status:    1,
		Visible:   1,
	}
	if err := db.Create(&operationLogMenu).Error; err != nil {
		fmt.Printf("创建操作日志菜单失败: %v\n", err)
		return
	}

	// 创建登录日志菜单
	loginLogMenu := models.SystemMenu{
		ParentID:  systemDir.ID,
		MenuName:  "登录日志",
		MenuType:  2, // 菜单
		Icon:      "ri:login-box-line",
		Path:      "/system/log/login",
		Component: "system/log/login",
		Perm:      "system:log:login:view",
		Sort:      6,
		Status:    1,
		Visible:   1,
	}
	if err := db.Create(&loginLogMenu).Error; err != nil {
		fmt.Printf("创建登录日志菜单失败: %v\n", err)
		return
	}

	// 查询 admin 角色 ID
	var adminRole models.SystemRole
	if err := db.Where("role_code = ?", "admin").First(&adminRole).Error; err != nil {
		fmt.Printf("未找到 admin 角色，跳过菜单权限分配: %v\n", err)
		return
	}

	// 给超级管理员角色分配所有菜单权限（首页使用前端静态路由）
	menus := []models.SystemMenu{systemDir, userMenu, roleMenu, menuMenu, settingMenu, operationLogMenu, loginLogMenu}
	for _, menu := range menus {
		roleMenu := models.SystemRoleMenu{
			RoleID: adminRole.ID,
			MenuID: menu.ID,
		}
		if err := db.Create(&roleMenu).Error; err != nil {
			fmt.Printf("分配菜单权限失败: %v\n", err)
			return
		}
	}

	fmt.Println("菜单数据初始化完成")
}

// initSystemSetting 初始化系统设置
func initSystemSetting() {
	db := global.DB

	// 检查是否已存在设置
	var count int64
	db.Model(&models.SystemSetting{}).Count(&count)
	if count > 0 {
		fmt.Println("系统设置已存在，跳过初始化")
		return
	}

	// 创建默认设置（默认关闭日志）
	setting := models.SystemSetting{
		EnableOperationLog: 2, // 默认关闭
		EnableLoginLog:     2, // 默认关闭
	}
	if err := db.Create(&setting).Error; err != nil {
		fmt.Printf("创建系统设置失败: %v\n", err)
		return
	}

	fmt.Println("系统设置初始化完成（默认关闭操作日志和登录日志）")
}

// initCasbinPolicy 初始化Casbin策略
func initCasbinPolicy(roleID uint) {
	// 创建Casbin适配器
	adapter, err := gormadapter.NewAdapterByDB(global.DB)
	if err != nil {
		fmt.Printf("创建Casbin适配器失败: %v\n", err)
		return
	}

	// 加载模型配置（优先使用配置文件，如果不存在则使用内联配置）
	var m model.Model
	modelPath := global.Config.Casbin.ModelPath
	if modelPath != "" {
		if _, err := os.Stat(modelPath); err == nil {
			// 配置文件存在，从文件加载
			m, err = model.NewModelFromFile(modelPath)
			if err != nil {
				fmt.Printf("从文件加载Casbin模型失败: %v\n", err)
				return
			}
			fmt.Printf("Casbin模型已从文件加载: %s\n", modelPath)
		} else {
			// 配置文件不存在，使用内联配置
			m, err = model.NewModelFromString(global.Config.System.CasbinConfig)
			if err != nil {
				fmt.Printf("加载Casbin模型失败: %v\n", err)
				return
			}
			fmt.Println("Casbin模型已使用内联配置加载")
		}
	} else {
		// 未配置模型路径，使用内联配置
		m, err = model.NewModelFromString(global.Config.System.CasbinConfig)
		if err != nil {
			fmt.Printf("加载Casbin模型失败: %v\n", err)
			return
		}
		fmt.Println("Casbin模型已使用内联配置加载")
	}

	// 创建Enforcer
	e, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		fmt.Printf("创建Casbin Enforcer失败: %v\n", err)
		return
	}

	// 添加策略：超级管理员拥有所有权限
	adminRoleKey := fmt.Sprintf("role_%d", roleID)
	e.AddPolicy(adminRoleKey, "*", "*")
	fmt.Printf("已为超级管理员(%s)添加通配符权限\n", adminRoleKey)

	// 为所有角色添加基础权限
	syncAllRolesCasbinPolicy(e)

	fmt.Println("Casbin权限初始化完成")
}

// initCasbinPolicyForAllRoles 为所有角色初始化Casbin策略（用于数据已存在的情况）
func initCasbinPolicyForAllRoles() {
	// 创建Casbin适配器
	adapter, err := gormadapter.NewAdapterByDB(global.DB)
	if err != nil {
		fmt.Printf("创建Casbin适配器失败: %v\n", err)
		return
	}

	// 加载模型配置
	var m model.Model
	modelPath := global.Config.Casbin.ModelPath
	if modelPath != "" {
		if _, err := os.Stat(modelPath); err == nil {
			m, err = model.NewModelFromFile(modelPath)
			if err != nil {
				fmt.Printf("加载Casbin模型失败: %v\n", err)
				return
			}
		} else {
			m, err = model.NewModelFromString(global.Config.System.CasbinConfig)
			if err != nil {
				fmt.Printf("加载Casbin模型失败: %v\n", err)
				return
			}
		}
	} else {
		m, err = model.NewModelFromString(global.Config.System.CasbinConfig)
		if err != nil {
			fmt.Printf("加载Casbin模型失败: %v\n", err)
			return
		}
	}

	// 创建Enforcer
	e, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		fmt.Printf("创建Casbin Enforcer失败: %v\n", err)
		return
	}

	// 为所有角色同步策略
	syncAllRolesCasbinPolicy(e)
}

// syncAllRolesCasbinPolicy 为所有角色同步Casbin策略
func syncAllRolesCasbinPolicy(e *casbin.Enforcer) {
	// 查询所有角色
	var roles []models.SystemRole
	if err := global.DB.Find(&roles).Error; err != nil {
		fmt.Printf("查询角色列表失败: %v\n", err)
		return
	}

	for _, role := range roles {
		roleKey := fmt.Sprintf("role_%d", role.ID)
		
		// 查询该角色的菜单权限
		var menuIDs []uint
		global.DB.Model(&models.SystemRoleMenu{}).Where("role_id = ?", role.ID).Pluck("menu_id", &menuIDs)
		
		// 所有角色都允许访问菜单路由接口（注意路径带 /api 前缀）
		_, _ = e.AddPolicy(roleKey, "/api/v1/system/routes", "GET")
		
		// 所有角色都允许访问个人中心相关接口
		e.AddPolicy(roleKey, "/api/v1/system/users/info", "GET")
		e.AddPolicy(roleKey, "/api/v1/system/users/profile", "PUT")
		e.AddPolicy(roleKey, "/api/v1/system/users/password", "PUT")
		
		// 根据菜单权限添加对应的API权限（路径带 /api 前缀）
		for _, menuID := range menuIDs {
			var menu models.SystemMenu
			if err := global.DB.First(&menu, menuID).Error; err != nil {
				continue
			}
			
			switch menu.MenuType {
			case 2: // 菜单
				apiPath := mapMenuPathToAPI(menu.Path)
				if apiPath != "" {
					e.AddPolicy(roleKey, "/api"+apiPath, "GET")
				}
			case 3: // 按钮
				if menu.Perm != "" {
					apiPath := mapPermToAPI(menu.Perm)
					if apiPath != "" {
						method := mapPermToMethod(menu.Perm)
						e.AddPolicy(roleKey, "/api"+apiPath, method)
					}
				}
			}
		}
		
		fmt.Printf("已为角色 %s(%s) 同步 %d 个菜单的权限\n", role.RoleName, roleKey, len(menuIDs))
	}
	
	// 保存策略到数据库
	_ = e.SavePolicy()
}

// mapMenuPathToAPI 将菜单路径映射到API路径
func mapMenuPathToAPI(menuPath string) string {
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

// mapPermToAPI 将权限标识映射到API路径
func mapPermToAPI(perm string) string {
	parts := strings.Split(perm, ":")
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

// mapPermToMethod 将权限标识映射到HTTP方法
func mapPermToMethod(perm string) string {
	if strings.Contains(perm, ":add") {
		return "POST"
	}
	if strings.Contains(perm, ":edit") {
		return "PUT"
	}
	if strings.Contains(perm, ":delete") {
		return "DELETE"
	}
	return "GET"
}
