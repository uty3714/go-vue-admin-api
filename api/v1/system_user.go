package v1

import (
	"go-vue-admin/models"
	"go-vue-admin/models/res"
	"go-vue-admin/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type SystemUserApi struct{}

// ==================== 认证相关 ====================

// Login 用户登录
func (a *SystemUserApi) Login(c *gin.Context) {
	var req models.SystemUserLoginReq
	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		res.ValidationError(c, err.Error())
		return
	}

	resp, errCode := systemUserService.Login(&req, c.ClientIP())
	if errCode != res.SuccessCode {
		res.Fail(c, errCode)
		return
	}

	res.Success(c, resp)
}

// Logout 用户登出
func (a *SystemUserApi) Logout(c *gin.Context) {
	// TODO: 可以将token加入黑名单
	res.Success(c, nil)
}

// GetUserInfo 获取当前用户信息
func (a *SystemUserApi) GetUserInfo(c *gin.Context) {
	userId, _ := c.Get("userId")

	user, err := systemUserService.GetUserInfo(userId.(uint))
	if err != nil {
		res.Fail(c, res.ErrorCodeUserNotExist)
		return
	}

	res.Success(c, user)
}

// GetAsyncRoutes 获取当前用户的动态路由菜单
func (a *SystemUserApi) GetAsyncRoutes(c *gin.Context) {
	// 从JWT中获取用户ID
	userId, exists := c.Get("userId")
	if !exists {
		res.Fail(c, res.ErrorCodeUnauthorized)
		return
	}

	routes, errCode := systemUserService.GetAsyncRoutes(userId.(uint))
	if errCode != res.SuccessCode {
		res.Fail(c, errCode)
		return
	}

	res.Success(c, routes)
}

// ==================== 用户管理 ====================

// GetUserList 获取用户列表
func (a *SystemUserApi) GetUserList(c *gin.Context) {
	page := util.StringToInt(c.DefaultQuery("page", "1"))
	pageSize := util.StringToInt(c.DefaultQuery("pageSize", "10"))
	keyword := c.Query("keyword")
	status := c.Query("status")

	users, total, err := systemUserService.GetUserList(page, pageSize, keyword, status)
	if err != nil {
		res.Error(c, err)
		return
	}

	res.PageSuccess(c, users, total, page, pageSize)
}

// CreateUser 创建用户
func (a *SystemUserApi) CreateUser(c *gin.Context) {
	var req models.SystemUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		res.ValidationError(c, err.Error())
		return
	}

	// 检查用户名是否已存在
	if systemUserService.CheckUserExist(req.Username) {
		res.Fail(c, res.ErrorCodeUserExist)
		return
	}

	id, err := systemUserService.CreateUser(&req)
	if err != nil {
		res.Error(c, err)
		return
	}

	res.Success(c, id)
}

// UpdateUser 更新用户
func (a *SystemUserApi) UpdateUser(c *gin.Context) {
	var req models.SystemUserUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		res.ValidationError(c, err.Error())
		return
	}

	if req.ID == 0 {
		res.Fail(c, res.ErrorCodeParamInvalid)
		return
	}

	// 检查用户是否存在
	if _, err := systemUserService.GetUserByID(req.ID); err != nil {
		res.Fail(c, res.ErrorCodeUserNotExist)
		return
	}

	// 如果要修改用户名，检查是否与其他用户冲突
	if req.Username != "" && systemUserService.CheckUserExistExceptID(req.Username, req.ID) {
		res.Fail(c, res.ErrorCodeUserExist)
		return
	}

	if err := systemUserService.UpdateUser(&req); err != nil {
		res.Error(c, err)
		return
	}

	res.Success(c, nil)
}

// DeleteUser 删除用户
func (a *SystemUserApi) DeleteUser(c *gin.Context) {
	id := util.StringToUint(c.Param("id"))
	if id == 0 {
		res.Fail(c, res.ErrorCodeParamInvalid)
		return
	}

	if err := systemUserService.DeleteUser(id); err != nil {
		res.Error(c, err)
		return
	}

	res.Success(c, nil)
}
