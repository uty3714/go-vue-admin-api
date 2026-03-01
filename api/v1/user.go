package v1

import (
	"go-vue-admin/models"
	"go-vue-admin/models/res"
	"go-vue-admin/util"

	"github.com/gin-gonic/gin"
)

type UserApi struct{}

// ==================== 用户信息 ====================

// GetUserInfo 获取用户信息
func (a *UserApi) GetUserInfo(c *gin.Context) {
	userId, _ := c.Get("userId")

	user, err := userService.GetUserByID(userId.(uint))
	if err != nil {
		res.Fail(c, res.ErrorCodeUserNotExist)
		return
	}

	res.Success(c, user)
}

// UpdateUser 更新用户信息
func (a *UserApi) UpdateUser(c *gin.Context) {
	userId, _ := c.Get("userId")

	var req models.UserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		res.ValidationError(c, err.Error())
		return
	}

	if err := userService.UpdateUserInfo(userId.(uint), &req); err != nil {
		res.Error(c, err)
		return
	}

	res.Success(c, nil)
}

// ==================== 地址管理 ====================

// GetAddressList 获取地址列表
func (a *UserApi) GetAddressList(c *gin.Context) {
	userId, _ := c.Get("userId")

	addresses, err := userService.GetAddressList(userId.(uint))
	if err != nil {
		res.Error(c, err)
		return
	}

	res.Success(c, addresses)
}

// CreateAddress 创建地址
func (a *UserApi) CreateAddress(c *gin.Context) {
	userId, _ := c.Get("userId")

	var req models.UserAddressReq
	if err := c.ShouldBindJSON(&req); err != nil {
		res.ValidationError(c, err.Error())
		return
	}

	id, err := userService.CreateAddress(userId.(uint), &req)
	if err != nil {
		res.Error(c, err)
		return
	}

	res.Success(c, id)
}

// UpdateAddress 更新地址
func (a *UserApi) UpdateAddress(c *gin.Context) {
	userId, _ := c.Get("userId")
	id := util.StringToUint(c.Param("id"))
	if id == 0 {
		res.Fail(c, res.ErrorCodeParamInvalid)
		return
	}

	var req models.UserAddressReq
	if err := c.ShouldBindJSON(&req); err != nil {
		res.ValidationError(c, err.Error())
		return
	}

	if err := userService.UpdateAddress(id, userId.(uint), &req); err != nil {
		res.Error(c, err)
		return
	}

	res.Success(c, nil)
}

// DeleteAddress 删除地址
func (a *UserApi) DeleteAddress(c *gin.Context) {
	userId, _ := c.Get("userId")
	id := util.StringToUint(c.Param("id"))
	if id == 0 {
		res.Fail(c, res.ErrorCodeParamInvalid)
		return
	}

	if err := userService.DeleteAddress(id, userId.(uint)); err != nil {
		res.Error(c, err)
		return
	}

	res.Success(c, nil)
}

// ==================== 购物车管理 ====================

// GetCartList 获取购物车列表
func (a *UserApi) GetCartList(c *gin.Context) {
	userId, _ := c.Get("userId")

	carts, err := userService.GetCartList(userId.(uint))
	if err != nil {
		res.Error(c, err)
		return
	}

	res.Success(c, carts)
}

// AddCart 添加购物车
func (a *UserApi) AddCart(c *gin.Context) {
	userId, _ := c.Get("userId")

	var req models.UserCartReq
	if err := c.ShouldBindJSON(&req); err != nil {
		res.ValidationError(c, err.Error())
		return
	}

	id, err := userService.AddOrUpdateCart(userId.(uint), &req)
	if err != nil {
		res.Error(c, err)
		return
	}

	res.Success(c, id)
}

// UpdateCart 更新购物车
func (a *UserApi) UpdateCart(c *gin.Context) {
	userId, _ := c.Get("userId")
	id := util.StringToUint(c.Param("id"))
	if id == 0 {
		res.Fail(c, res.ErrorCodeParamInvalid)
		return
	}

	var req models.UserCartUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		res.ValidationError(c, err.Error())
		return
	}

	if err := userService.UpdateCart(id, userId.(uint), &req); err != nil {
		res.Error(c, err)
		return
	}

	res.Success(c, nil)
}

// DeleteCart 删除购物车
func (a *UserApi) DeleteCart(c *gin.Context) {
	userId, _ := c.Get("userId")
	id := util.StringToUint(c.Param("id"))
	if id == 0 {
		res.Fail(c, res.ErrorCodeParamInvalid)
		return
	}

	if err := userService.DeleteCart(id, userId.(uint)); err != nil {
		res.Error(c, err)
		return
	}

	res.Success(c, nil)
}
