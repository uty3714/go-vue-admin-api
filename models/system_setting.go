package models

// SystemSetting 系统设置表
type SystemSetting struct {
	ID                uint      `gorm:"column:id;primarykey;comment:主键ID" json:"id"`
	CreatedAt         LocalTime `gorm:"column:created_at;type:datetime;not null;comment:创建时间" json:"createdAt"`
	UpdatedAt         LocalTime `gorm:"column:updated_at;type:datetime;not null;comment:更新时间" json:"updatedAt"`
	EnableOperationLog int      `gorm:"column:enable_operation_log;type:tinyint;default:2;comment:是否开启操作日志 1开启 2关闭" json:"enableOperationLog"`
	EnableLoginLog    int       `gorm:"column:enable_login_log;type:tinyint;default:2;comment:是否开启登录日志 1开启 2关闭" json:"enableLoginLog"`
}

func (SystemSetting) TableName() string {
	return "system_setting"
}
