package models

// UserIntegral 用户积分记录表
type UserIntegral struct {
	Model
	UserID      uint   `gorm:"not null;index;comment:用户ID" json:"userId"`
	Type        int    `gorm:"type:tinyint;not null;comment:类型 1获得 2消耗" json:"type"`
	Integral    int    `gorm:"type:int;not null;comment:积分数量" json:"integral"`
	Balance     int    `gorm:"type:int;not null;comment:积分余额" json:"balance"`
	Source      string `gorm:"type:varchar(255);comment:来源" json:"source"`
	SourceID    uint   `gorm:"comment:来源ID" json:"sourceId"`
	Description string `gorm:"type:varchar(255);comment:描述" json:"description"`
}

func (UserIntegral) TableName() string {
	return "user_integral"
}
