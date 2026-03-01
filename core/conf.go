package core

import (
	"go-vue-admin/conf"
	"go-vue-admin/global"

	"github.com/sirupsen/logrus"
)

func InitConf(path string) {
	serverConfig, err := conf.InitConfig(path)
	if err != nil {
		logrus.Fatalf("加载配置文件失败: %v", err)
	}
	global.Config = serverConfig
	logrus.Info("配置文件加载成功")
}
