package v1

import (
	"errors"
	"go-vue-admin/global"
	"go-vue-admin/models"
	"time"
)

type SystemSettingService struct{}

// GetSetting 获取系统设置（如果不存在则创建默认设置）
func (s *SystemSettingService) GetSetting() (*models.SystemSetting, error) {
	var setting models.SystemSetting
	err := global.DB.First(&setting).Error
	if err != nil {
		// 如果没有记录，创建默认设置（默认关闭）
		setting = models.SystemSetting{
			EnableOperationLog: 2, // 默认关闭
			EnableLoginLog:     2, // 默认关闭
		}
		if createErr := global.DB.Create(&setting).Error; createErr != nil {
			return nil, createErr
		}
	}
	return &setting, nil
}

// UpdateSetting 更新系统设置
func (s *SystemSettingService) UpdateSetting(setting *models.SystemSetting) error {
	if setting == nil {
		return errors.New("设置信息不能为空")
	}

	// 获取现有设置或创建新设置
	var existing models.SystemSetting
	err := global.DB.First(&existing).Error
	if err != nil {
		// 不存在则创建
		setting.ID = 0 // 确保是新记录
		return global.DB.Create(setting).Error
	}

	// 更新设置
	setting.ID = existing.ID
	setting.CreatedAt = existing.CreatedAt
	return global.DB.Model(&existing).Updates(map[string]interface{}{
		"enable_operation_log": setting.EnableOperationLog,
		"enable_login_log":     setting.EnableLoginLog,
		"updated_at":           time.Now(),
	}).Error
}

// IsOperationLogEnabled 检查操作日志是否开启
func (s *SystemSettingService) IsOperationLogEnabled() bool {
	setting, err := s.GetSetting()
	if err != nil {
		return false // 默认关闭
	}
	return setting.EnableOperationLog == 1
}

// IsLoginLogEnabled 检查登录日志是否开启
func (s *SystemSettingService) IsLoginLogEnabled() bool {
	setting, err := s.GetSetting()
	if err != nil {
		return false // 默认关闭
	}
	return setting.EnableLoginLog == 1
}
