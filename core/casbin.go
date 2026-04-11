package core

import (
	"fmt"
	"go-vue-admin/global"
	"go-vue-admin/models"
	"os"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

// InitCasbin 初始化Casbin权限管理
func InitCasbin() *casbin.Enforcer {
	// 创建Casbin适配器
	adapter, err := gormadapter.NewAdapterByDB(global.DB)
	if err != nil {
		global.Log.Errorf("创建Casbin适配器失败: %v", err)
		return nil
	}

	// 加载模型配置
	var m model.Model
	modelPath := global.Config.Casbin.ModelPath
	if modelPath != "" {
		if _, err := os.Stat(modelPath); err == nil {
			m, err = model.NewModelFromFile(modelPath)
			if err != nil {
				global.Log.Errorf("从文件加载Casbin模型失败: %v", err)
				return nil
			}
		} else {
			m, err = model.NewModelFromString(global.Config.System.CasbinConfig)
			if err != nil {
				global.Log.Errorf("加载Casbin模型失败: %v", err)
				return nil
			}
		}
	} else {
		m, err = model.NewModelFromString(global.Config.System.CasbinConfig)
		if err != nil {
			global.Log.Errorf("加载Casbin模型失败: %v", err)
			return nil
		}
	}

	// 创建Enforcer
	e, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		global.Log.Errorf("创建Casbin Enforcer失败: %v", err)
		return nil
	}

	// 从数据库加载策略
	if err := e.LoadPolicy(); err != nil {
		global.Log.Errorf("加载Casbin策略失败: %v", err)
		return nil
	}

	// 检查是否有策略数据，如果没有则初始化
	policies := e.GetPolicy()
	if len(policies) == 0 {
		global.Log.Info("未发现Casbin策略，开始初始化...")
		if err := initCasbinPolicies(e); err != nil {
			global.Log.Errorf("初始化Casbin策略失败: %v", err)
		}
	} else {
		global.Log.Infof("已从数据库加载 %d 条策略", len(policies))
	}

	global.Log.Info("Casbin权限管理初始化成功")
	return e
}

// initCasbinPolicies 初始化所有角色的策略
func initCasbinPolicies(e *casbin.Enforcer) error {
	// 查询所有角色
	var roles []models.SystemRole
	if err := global.DB.Find(&roles).Error; err != nil {
		return fmt.Errorf("查询角色列表失败: %v", err)
	}

	for _, role := range roles {
		roleKey := fmt.Sprintf("role_%d", role.ID)

		// 查询该角色的菜单权限
		var menuIDs []uint
		global.DB.Model(&models.SystemRoleMenu{}).Where("role_id = ?", role.ID).Pluck("menu_id", &menuIDs)

		global.Log.Infof("[Casbin] 初始化角色 %s 的策略，菜单数: %d", roleKey, len(menuIDs))

		// 所有角色都允许访问菜单路由接口
		e.AddPolicy(roleKey, "/api/v1/system/routes", "GET")

		// 所有角色都允许访问个人中心相关接口
		e.AddPolicy(roleKey, "/api/v1/system/users/info", "GET")
		e.AddPolicy(roleKey, "/api/v1/system/users/profile", "PUT")
		e.AddPolicy(roleKey, "/api/v1/system/users/password", "PUT")

		// 根据菜单权限添加对应的API权限
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

		global.Log.Infof("[Casbin] 角色 %s 策略初始化完成", roleKey)
	}

	// 保存策略到数据库
	if err := e.SavePolicy(); err != nil {
		return fmt.Errorf("保存策略失败: %v", err)
	}

	global.Log.Infof("[Casbin] 已为 %d 个角色初始化策略并保存到数据库", len(roles))
	return nil
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
