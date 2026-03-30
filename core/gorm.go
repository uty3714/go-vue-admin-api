package core

import (
	"log"
	"os"
	"time"
	"go-vue-admin/global"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func InitGorm() *gorm.DB {
	if global.Config.System.DbType == "mysql" {
		return InitGormMysql()
	}
	return nil
}

func InitGormMysql() *gorm.DB {
	m := global.Config.Mysql

	if m.DbName == "" {
		return nil
	}

	dsn := m.Dsn()
	
	var logMode logger.Interface
	if m.LogMode == "info" {
		logMode = logger.Default.LogMode(logger.Info)
	} else if m.LogMode == "warn" {
		logMode = logger.Default.LogMode(logger.Warn)
	} else if m.LogMode == "error" {
		logMode = logger.Default.LogMode(logger.Error)
	} else {
		logMode = logger.Default.LogMode(logger.Silent)
	}
	
	if m.LogZap {
		logMode = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: true,
				ParameterizedQueries:      false,
				Colorful:                  true,
			},
		)
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		Logger:         logMode,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		global.Log.Errorf("连接mysql数据库失败: %v", err)
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		global.Log.Errorf("获取sqlDB失败: %v", err)
		return nil
	}
	
	sqlDB.SetMaxIdleConns(m.MaxIdleConns)
	sqlDB.SetMaxOpenConns(m.MaxOpenConns)
	// 添加连接最大生命周期配置，防止连接泄漏和超时问题
	sqlDB.SetConnMaxLifetime(time.Hour)
	// 添加连接最大空闲时间
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	return db
}
