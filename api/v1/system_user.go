package v1

import (
	"go-vue-admin/middleware"
	"go-vue-admin/models"
	"go-vue-admin/models/res"
	"go-vue-admin/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type SystemUserApi struct{}

// ==================== 认证相关 ====================

// Login
// @Tags 系统管理-认证
// @Summary 用户登录
// @Description 用户登录接口，返回JWT token
// @Accept json
// @Produce json
// @Param data body models.SystemUserLoginReq true "登录参数"
// @Success 200 {object} res.Response{data=models.SystemUserLoginRes} "登录成功"
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 401 {object} res.Response "登录失败，用户名或密码错误"
// @Router /api/v1/system/login [post]
func (a *SystemUserApi) Login(c *gin.Context) {
	var req models.SystemUserLoginReq
	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		res.ValidationError(c, err.Error())
		return
	}

	// 获取客户端信息
	userAgent := c.Request.UserAgent()
	resp, errCode := systemUserService.Login(&req, c.ClientIP(), userAgent)
	if errCode != res.SuccessCode {
		res.Fail(c, errCode)
		return
	}

	res.Success(c, resp)
}

// Logout
// @Tags 系统管理-认证
// @Summary 用户登出
// @Description 用户登出，token将被加入黑名单
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} res.Response "登出成功"
// @Router /api/v1/system/logout [post]
func (a *SystemUserApi) Logout(c *gin.Context) {
	middleware.LogoutHandler(c)
}

// GetUserInfo
// @Tags 系统管理-认证
// @Summary 获取当前用户信息
// @Description 获取当前登录用户的详细信息
// @Produce json
// @Security BearerAuth
// @Success 200 {object} res.Response{data=models.SystemUser} "成功"
// @Failure 401 {object} res.Response "未登录或token过期"
// @Router /api/v1/system/user/info [get]
func (a *SystemUserApi) GetUserInfo(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		res.Fail(c, res.ErrorCodeUnauthorized)
		return
	}
	
	// 安全的类型断言
	uid, ok := userId.(uint)
	if !ok {
		res.Fail(c, res.ErrorCodeUnauthorized)
		return
	}

	user, err := systemUserService.GetUserInfo(uid)
	if err != nil {
		res.Fail(c, res.ErrorCodeUserNotExist)
		return
	}

	res.Success(c, user)
}

// GetAsyncRoutes
// @Tags 系统管理-认证
// @Summary 获取动态路由
// @Description 获取当前登录用户的动态路由菜单
// @Produce json
// @Security BearerAuth
// @Success 200 {object} res.Response{data=[]map[string]interface{}} "成功"
// @Failure 401 {object} res.Response "未登录或token过期"
// @Router /api/v1/system/routes [get]
func (a *SystemUserApi) GetAsyncRoutes(c *gin.Context) {
	// 从JWT中获取用户ID
	userId, exists := c.Get("userId")
	if !exists {
		res.Fail(c, res.ErrorCodeUnauthorized)
		return
	}

	// 安全的类型断言
	uid, ok := userId.(uint)
	if !ok {
		res.Fail(c, res.ErrorCodeUnauthorized)
		return
	}

	routes, errCode := systemUserService.GetAsyncRoutes(uid)
	if errCode != res.SuccessCode {
		res.Fail(c, errCode)
		return
	}

	res.Success(c, routes)
}

// ==================== 用户管理 ====================

// GetUserList
// @Tags 系统管理-用户
// @Summary 获取用户列表
// @Description 分页获取系统用户列表，支持关键词搜索和状态筛选
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码，默认1"
// @Param pageSize query int false "每页数量，默认10，最大100"
// @Param keyword query string false "关键词（用户名/昵称/手机号）"
// @Param status query string false "状态：1启用 2禁用"
// @Success 200 {object} res.Response{data=res.PageResult{list=[]models.SystemUser}} "成功"
// @Failure 401 {object} res.Response "未登录或token过期"
// @Router /api/v1/system/user/list [get]
func (a *SystemUserApi) GetUserList(c *gin.Context) {
	page := util.StringToInt(c.DefaultQuery("page", "1"))
	pageSize := util.StringToInt(c.DefaultQuery("pageSize", "10"))
	
	// 限制分页大小，防止性能问题
	const maxPageSize = 100
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if page < 1 {
		page = 1
	}
	
	keyword := c.Query("keyword")
	status := c.Query("status")

	users, total, err := systemUserService.GetUserList(page, pageSize, keyword, status)
	if err != nil {
		res.Error(c, err)
		return
	}

	res.PageSuccess(c, users, total, page, pageSize)
}

// CreateUser
// @Tags 系统管理-用户
// @Summary 创建用户
// @Description 创建新用户，用户名不能重复
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body models.SystemUserReq true "用户数据"
// @Success 200 {object} res.Response{data=uint} "创建成功，返回用户ID"
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 401 {object} res.Response "未登录或token过期"
// @Router /api/v1/system/user/create [post]
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

// UpdateUser
// @Tags 系统管理-用户
// @Summary 更新用户
// @Description 更新用户信息，用户名不能与其他用户重复。管理员不能通过此接口修改密码，请使用密码重置功能。
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body models.SystemUserUpdateReq true "用户数据"
// @Success 200 {object} res.Response "更新成功"
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 401 {object} res.Response "未登录或token过期"
// @Failure 404 {object} res.Response "用户不存在"
// @Router /api/v1/system/user/update [put]
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

// DeleteUser
// @Tags 系统管理-用户
// @Summary 删除用户
// @Description 根据ID删除用户（软删除）
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Success 200 {object} res.Response "删除成功"
// @Failure 401 {object} res.Response "未登录或token过期"
// @Failure 404 {object} res.Response "用户不存在"
// @Router /api/v1/system/user/delete/{id} [delete]
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

// UpdateCurrentUser
// @Tags 系统管理-用户
// @Summary 更新当前用户信息
// @Description 当前登录用户更新自己的个人信息（昵称、头像、手机号、邮箱）
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body models.SystemUserProfileReq true "个人信息"
// @Success 200 {object} res.Response "更新成功"
// @Failure 401 {object} res.Response "未登录或token过期"
// @Router /api/v1/system/user/profile [put]
func (a *SystemUserApi) UpdateCurrentUser(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		res.Fail(c, res.ErrorCodeUnauthorized)
		return
	}
	
	// 安全的类型断言
	uid, ok := userId.(uint)
	if !ok {
		res.Fail(c, res.ErrorCodeUnauthorized)
		return
	}

	var req models.SystemUserProfileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		res.ValidationError(c, err.Error())
		return
	}

	if err := systemUserService.UpdateCurrentUser(uid, &req); err != nil {
		res.Error(c, err)
		return
	}

	res.Success(c, nil)
}

// UpdateCurrentUserPassword
// @Tags 系统管理-用户
// @Summary 修改当前用户密码
// @Description 当前登录用户修改自己的密码，需要验证原密码
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body models.SystemUserPasswordReq true "密码数据"
// @Success 200 {object} res.Response "修改成功"
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 401 {object} res.Response "未登录或token过期/原密码错误"
// @Router /api/v1/system/user/password [put]
func (a *SystemUserApi) UpdateCurrentUserPassword(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		res.Fail(c, res.ErrorCodeUnauthorized)
		return
	}
	
	// 安全的类型断言
	uid, ok := userId.(uint)
	if !ok {
		res.Fail(c, res.ErrorCodeUnauthorized)
		return
	}

	var req models.SystemUserPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		res.ValidationError(c, err.Error())
		return
	}

	if err := systemUserService.UpdateCurrentUserPassword(uid, req.OldPassword, req.NewPassword); err != nil {
		res.FailWithMessage(c, res.ErrorCodeBusinessError, err.Error())
		return
	}

	res.Success(c, nil)
}
