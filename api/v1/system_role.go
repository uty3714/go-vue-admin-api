package v1

import (
	"errors"
	"go-vue-admin/models"
	"go-vue-admin/models/res"
	"go-vue-admin/services/v1"
	"go-vue-admin/util"

	"github.com/gin-gonic/gin"
)

type SystemRoleApi struct{}

// SetRoleMenusReq 设置角色菜单权限请求（复用 service 层的定义）
type SetRoleMenusReq = v1.SetRoleMenusReq

// ==================== 角色管理 ====================

// GetRoleList 获取角色列表
func (a *SystemRoleApi) GetRoleList(c *gin.Context) {
	page := util.StringToInt(c.DefaultQuery("page", "1"))
	pageSize := util.StringToInt(c.DefaultQuery("pageSize", "10"))
	keyword := c.Query("keyword")

	roles, total, err := systemRoleService.GetRoleList(page, pageSize, keyword)
	if err != nil {
		res.Error(c, err)
		return
	}

	res.PageSuccess(c, roles, total, page, pageSize)
}

// GetRoleDetail 获取角色详情
func (a *SystemRoleApi) GetRoleDetail(c *gin.Context) {
	id := util.StringToUint(c.Param("id"))
	if id == 0 {
		res.Fail(c, res.ErrorCodeParamInvalid)
		return
	}

	role, err := systemRoleService.GetRoleByID(id)
	if err != nil {
		res.Fail(c, res.ErrorCodeNotFound)
		return
	}

	res.Success(c, role)
}

// CreateRole 创建角色
func (a *SystemRoleApi) CreateRole(c *gin.Context) {
	var req models.SystemRole
	if err := c.ShouldBindJSON(&req); err != nil {
		res.ValidationError(c, err.Error())
		return
	}

	// 检查角色代码是否已存在
	if systemRoleService.CheckRoleCodeExist(req.RoleCode) {
		res.FailWithMessage(c, res.ErrorCodeBusinessError, "角色代码已存在")
		return
	}

	id, err := systemRoleService.CreateRole(&req)
	if err != nil {
		res.Error(c, err)
		return
	}

	res.Success(c, id)
}

// UpdateRole 更新角色
func (a *SystemRoleApi) UpdateRole(c *gin.Context) {
	var req models.SystemRole
	if err := c.ShouldBindJSON(&req); err != nil {
		res.ValidationError(c, err.Error())
		return
	}

	if req.ID == 0 {
		res.Fail(c, res.ErrorCodeParamInvalid)
		return
	}

	if err := systemRoleService.UpdateRole(&req); err != nil {
		res.Error(c, err)
		return
	}

	res.Success(c, nil)
}

// DeleteRole 删除角色
func (a *SystemRoleApi) DeleteRole(c *gin.Context) {
	id := util.StringToUint(c.Param("id"))
	if id == 0 {
		res.Fail(c, res.ErrorCodeParamInvalid)
		return
	}

	if err := systemRoleService.DeleteRole(id); err != nil {
		if errors.Is(err, v1.ErrRoleHasUsers) {
			res.FailWithMessage(c, res.ErrorCodeBusinessError, err.Error())
			return
		}
		res.Error(c, err)
		return
	}

	res.Success(c, nil)
}

// GetRoleMenus 获取角色的菜单权限
func (a *SystemRoleApi) GetRoleMenus(c *gin.Context) {
	id := util.StringToUint(c.Param("id"))
	if id == 0 {
		res.Fail(c, res.ErrorCodeParamInvalid)
		return
	}

	menuIDs, err := systemRoleService.GetRoleMenus(id)
	if err != nil {
		res.Error(c, err)
		return
	}

	res.Success(c, menuIDs)
}

// SetRoleMenus 设置角色的菜单权限
func (a *SystemRoleApi) SetRoleMenus(c *gin.Context) {
	var req SetRoleMenusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		res.ValidationError(c, err.Error())
		return
	}

	if err := systemRoleService.SetRoleMenus(&req); err != nil {
		res.Error(c, err)
		return
	}

	res.Success(c, nil)
}

// ==================== 菜单管理 ====================

// GetMenuList 获取菜单列表
func (a *SystemRoleApi) GetMenuList(c *gin.Context) {
	menus, err := systemRoleService.GetMenuList()
	if err != nil {
		res.Error(c, err)
		return
	}
	res.Success(c, menus)
}

// GetMenuTree 获取菜单树
func (a *SystemRoleApi) GetMenuTree(c *gin.Context) {
	menuTree, err := systemRoleService.GetMenuTree()
	if err != nil {
		res.Error(c, err)
		return
	}
	res.Success(c, menuTree)
}

// CreateMenu 创建菜单
func (a *SystemRoleApi) CreateMenu(c *gin.Context) {
	var req models.SystemMenu
	if err := c.ShouldBindJSON(&req); err != nil {
		res.ValidationError(c, err.Error())
		return
	}

	id, err := systemRoleService.CreateMenu(&req)
	if err != nil {
		res.Error(c, err)
		return
	}

	res.Success(c, id)
}

// UpdateMenu 更新菜单
func (a *SystemRoleApi) UpdateMenu(c *gin.Context) {
	var req models.SystemMenu
	if err := c.ShouldBindJSON(&req); err != nil {
		res.ValidationError(c, err.Error())
		return
	}

	if req.ID == 0 {
		res.Fail(c, res.ErrorCodeParamInvalid)
		return
	}

	if err := systemRoleService.UpdateMenu(&req); err != nil {
		res.Error(c, err)
		return
	}

	res.Success(c, nil)
}

// DeleteMenu 删除菜单
func (a *SystemRoleApi) DeleteMenu(c *gin.Context) {
	id := util.StringToUint(c.Param("id"))
	if id == 0 {
		res.Fail(c, res.ErrorCodeParamInvalid)
		return
	}

	if err := systemRoleService.DeleteMenu(id); err != nil {
		if errors.Is(err, v1.ErrMenuHasChildren) {
			res.FailWithMessage(c, res.ErrorCodeBusinessError, err.Error())
			return
		}
		res.Error(c, err)
		return
	}

	res.Success(c, nil)
}
