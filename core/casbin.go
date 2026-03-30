package core

import (
	"go-vue-admin/global"
	"os"

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

	global.Log.Info("Casbin权限管理初始化成功")
	return e
}
