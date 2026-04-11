package v1

import (
	"go-vue-admin/models"
	"go-vue-admin/models/res"

	"github.com/gin-gonic/gin"
)

type SystemSettingApi struct{}

// UpdateSettingReq 更新系统设置请求
type UpdateSettingReq struct {
	EnableOperationLog int `json:"enableOperationLog" binding:"oneof=1 2"`
	EnableLoginLog     int `json:"enableLoginLog" binding:"oneof=1 2"`
}

// GetSystemSetting
// @Tags 系统管理-系统设置
// @Summary 获取系统设置
// @Description 获取系统设置信息（操作日志、登录日志开关状态）
// @Produce json
// @Security BearerAuth
// @Success 200 {object} res.Response{data=models.SystemSetting} "成功"
// @Failure 401 {object} res.Response "未登录或token过期"
// @Router /api/v1/system/settings [get]
func (a *SystemSettingApi) GetSystemSetting(c *gin.Context) {
	setting, err := systemSettingService.GetSetting()
	if err != nil {
		res.Error(c, err)
		return
	}
	res.Success(c, setting)
}

// UpdateSystemSetting
// @Tags 系统管理-系统设置
// @Summary 更新系统设置
// @Description 更新系统设置（开启/关闭操作日志、登录日志）
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body UpdateSettingReq true "设置信息"
// @Success 200 {object} res.Response "成功"
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 401 {object} res.Response "未登录或token过期"
// @Router /api/v1/system/settings [put]
func (a *SystemSettingApi) UpdateSystemSetting(c *gin.Context) {
	var req UpdateSettingReq
	if err := c.ShouldBindJSON(&req); err != nil {
		res.Error(c, err)
		return
	}

	setting := &models.SystemSetting{
		EnableOperationLog: req.EnableOperationLog,
		EnableLoginLog:     req.EnableLoginLog,
	}

	if err := systemSettingService.UpdateSetting(setting); err != nil {
		res.Error(c, err)
		return
	}

	res.Success(c, nil)
}

// IsOperationLogEnabled 检查操作日志是否开启（内部使用）
func IsOperationLogEnabled() bool {
	return systemSettingService.IsOperationLogEnabled()
}

// IsLoginLogEnabled 检查登录日志是否开启（内部使用）
func IsLoginLogEnabled() bool {
	return systemSettingService.IsLoginLogEnabled()
}

// GetSettingForMenu 获取设置用于菜单显示控制（内部使用）
func GetSettingForMenu() (*models.SystemSetting, error) {
	return systemSettingService.GetSetting()
}
