package models

import (
	"time"

	"gorm.io/gorm"
)

// Model 基础模型结构体
type Model struct {
	ID        uint           `gorm:"primarykey;comment:主键ID" json:"id"`
	CreatedAt time.Time      `gorm:"type:datetime;not null;comment:创建时间" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"type:datetime;not null;comment:更新时间" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间" json:"-"`
}

// BaseModel 基础模型（不含软删除）
type BaseModel struct {
	ID        uint      `gorm:"primarykey;comment:主键ID" json:"id"`
	CreatedAt time.Time `gorm:"type:datetime;not null;comment:创建时间" json:"createdAt"`
	UpdatedAt time.Time `gorm:"type:datetime;not null;comment:更新时间" json:"updatedAt"`
}
