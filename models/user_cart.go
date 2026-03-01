package models

// UserCart 用户购物车表
type UserCart struct {
	Model
	UserID    uint  `gorm:"not null;index;comment:用户ID" json:"userId"`
	ProductID uint  `gorm:"not null;index;comment:商品ID" json:"productId"`
	SkuID     uint  `gorm:"not null;index;comment:SKU ID" json:"skuId"`
	Quantity  int   `gorm:"type:int;default:1;comment:数量" json:"quantity"`
	Selected  int   `gorm:"type:tinyint;default:1;comment:是否选中 1是 0否" json:"selected"`
}

func (UserCart) TableName() string {
	return "user_cart"
}

// UserCartReq 购物车请求参数
type UserCartReq struct {
	ProductID uint `json:"productId" binding:"required"`
	SkuID     uint `json:"skuId" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
	Selected  int  `json:"selected"`
}

// UserCartUpdateReq 购物车更新请求参数
type UserCartUpdateReq struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
	Selected int `json:"selected"`
}
