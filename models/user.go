package models

// User 用户表（前台用户）
type User struct {
	Model
	OpenID      string `gorm:"type:varchar(128);uniqueIndex;comment:微信OpenID" json:"openId"`
	UnionID     string `gorm:"type:varchar(128);comment:微信UnionID" json:"unionId"`
	Nickname    string `gorm:"type:varchar(64);comment:昵称" json:"nickname"`
	Avatar      string `gorm:"type:varchar(255);comment:头像URL" json:"avatar"`
	Phone       string `gorm:"type:varchar(20);uniqueIndex;comment:手机号" json:"phone"`
	Email       string `gorm:"type:varchar(128);comment:邮箱" json:"email"`
	Status      int    `gorm:"type:tinyint;default:1;comment:状态 1启用 2禁用" json:"status"`
	Gender      int    `gorm:"type:tinyint;default:0;comment:性别 0未知 1男 2女" json:"gender"`
	Birthday    string `gorm:"type:date;comment:生日" json:"birthday"`
	Integral    int    `gorm:"type:int;default:0;comment:积分" json:"integral"`
	Balance     int64  `gorm:"type:bigint;default:0;comment:余额(分)" json:"balance"`
	LastLoginIP string `gorm:"type:varchar(128);comment:最后登录IP" json:"lastLoginIp"`
	LastLoginAt string `gorm:"type:datetime;comment:最后登录时间" json:"lastLoginAt"`
}

func (User) TableName() string {
	return "user"
}

// UserReq 用户请求参数
type UserReq struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Gender   int    `json:"gender"`
	Birthday string `json:"birthday"`
}
