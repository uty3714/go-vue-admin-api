package models

// UserAddress 用户地址表
type UserAddress struct {
	Model
	UserID       uint   `gorm:"not null;index;comment:用户ID" json:"userId"`
	UserName     string `gorm:"type:varchar(64);not null;comment:收货人姓名" json:"userName"`
	UserPhone    string `gorm:"type:varchar(20);not null;comment:收货人手机号" json:"userPhone"`
	Province     string `gorm:"type:varchar(64);not null;comment:省份" json:"province"`
	City         string `gorm:"type:varchar(64);not null;comment:城市" json:"city"`
	District     string `gorm:"type:varchar(64);not null;comment:区/县" json:"district"`
	Detail       string `gorm:"type:varchar(255);not null;comment:详细地址" json:"detail"`
	PostalCode   string `gorm:"type:varchar(20);comment:邮政编码" json:"postalCode"`
	IsDefault    int    `gorm:"type:tinyint;default:0;comment:是否默认 1是 0否" json:"isDefault"`
}

func (UserAddress) TableName() string {
	return "user_address"
}

// UserAddressReq 地址请求参数
type UserAddressReq struct {
	UserName  string `json:"userName" binding:"required"`
	UserPhone string `json:"userPhone" binding:"required"`
	Province  string `json:"province" binding:"required"`
	City      string `json:"city" binding:"required"`
	District  string `json:"district" binding:"required"`
	Detail    string `json:"detail" binding:"required"`
	PostalCode string `json:"postalCode"`
	IsDefault int    `json:"isDefault"`
}
