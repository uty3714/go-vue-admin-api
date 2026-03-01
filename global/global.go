package global

import (
	"go-vue-admin/conf"

	"github.com/casbin/casbin/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	Config *conf.Server
	DB     *gorm.DB
	Log    *logrus.Logger
	Casbin *casbin.Enforcer
)
