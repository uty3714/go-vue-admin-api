package models

// OperationLog 操作日志表
type OperationLog struct {
	ID            uint      `gorm:"primarykey;comment:主键ID" json:"id"`
	UserID        uint      `gorm:"index;comment:用户ID" json:"userId"`
	Username      string    `gorm:"type:varchar(64);comment:用户名" json:"username"`
	RoleName      string    `gorm:"type:varchar(64);comment:角色名" json:"roleName"`
	Method        string    `gorm:"type:varchar(20);comment:请求方法" json:"method"`
	Path          string    `gorm:"type:varchar(255);comment:请求路径" json:"path"`
	RequestData   string    `gorm:"type:text;comment:请求数据" json:"requestData"`
	ResponseData  string    `gorm:"type:text;comment:响应数据" json:"responseData"`
	Status        int       `gorm:"default:1;comment:状态 1成功 2失败" json:"status"`
	ErrorMessage  string    `gorm:"type:text;comment:错误信息" json:"errorMessage"`
	IP            string    `gorm:"type:varchar(128);comment:IP地址" json:"ip"`
	UserAgent     string    `gorm:"type:varchar(512);comment:用户代理" json:"userAgent"`
	OperationTime int       `gorm:"comment:操作耗时(ms)" json:"operationTime"`
	CreatedAt     LocalTime `gorm:"type:datetime;not null;comment:创建时间" json:"createdAt"`
}

func (OperationLog) TableName() string {
	return "system_operation_log"
}

// LoginLog 登录日志表
type LoginLog struct {
	ID        uint      `gorm:"primarykey;comment:主键ID" json:"id"`
	Username  string    `gorm:"type:varchar(64);index;comment:用户名" json:"username"`
	IP        string    `gorm:"type:varchar(128);comment:IP地址" json:"ip"`
	Location  string    `gorm:"type:varchar(255);comment:登录地点" json:"location"`
	Browser   string    `gorm:"type:varchar(255);comment:浏览器" json:"browser"`
	OS        string    `gorm:"type:varchar(255);comment:操作系统" json:"os"`
	Status    int       `gorm:"default:1;comment:状态 1成功 2失败" json:"status"`
	Message   string    `gorm:"type:varchar(255);comment:消息" json:"message"`
	CreatedAt LocalTime `gorm:"type:datetime;not null;comment:创建时间" json:"createdAt"`
}

func (LoginLog) TableName() string {
	return "system_login_log"
}
