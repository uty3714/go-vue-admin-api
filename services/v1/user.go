package v1

import (
	"go-vue-admin/global"
	"go-vue-admin/models"
)

type UserService struct{}

// ==================== 用户基础操作 ====================

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := global.DB.First(&user, id).Error
	return &user, err
}

// GetUserByOpenID 根据OpenID获取用户
func (s *UserService) GetUserByOpenID(openID string) (*models.User, error) {
	var user models.User
	err := global.DB.Where("open_id = ?", openID).First(&user).Error
	return &user, err
}

// GetUserByPhone 根据手机号获取用户
func (s *UserService) GetUserByPhone(phone string) (*models.User, error) {
	var user models.User
	err := global.DB.Where("phone = ?", phone).First(&user).Error
	return &user, err
}

// CheckUserExist 检查用户是否存在
func (s *UserService) CheckUserExist(openID string) bool {
	var count int64
	global.DB.Model(&models.User{}).Where("open_id = ?", openID).Count(&count)
	return count > 0
}

// CreateUser 创建用户
func (s *UserService) CreateUser(user *models.User) error {
	return global.DB.Create(user).Error
}

// UpdateUser 更新用户
func (s *UserService) UpdateUser(user *models.User) error {
	return global.DB.Save(user).Error
}

// UpdateUserInfo 更新用户信息（前台用户）
func (s *UserService) UpdateUserInfo(userID uint, req *models.UserReq) error {
	return global.DB.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"nickname": req.Nickname,
		"avatar":   req.Avatar,
		"phone":    req.Phone,
		"email":    req.Email,
		"gender":   req.Gender,
		"birthday": req.Birthday,
	}).Error
}

// ==================== 地址管理 ====================

// GetAddressList 获取地址列表
func (s *UserService) GetAddressList(userID uint) ([]models.UserAddress, error) {
	var addresses []models.UserAddress
	err := global.DB.Where("user_id = ?", userID).Find(&addresses).Error
	return addresses, err
}

// GetAddressByID 根据ID获取地址
func (s *UserService) GetAddressByID(id, userID uint) (*models.UserAddress, error) {
	var address models.UserAddress
	err := global.DB.Where("id = ? AND user_id = ?", id, userID).First(&address).Error
	return &address, err
}

// CreateAddress 创建地址
func (s *UserService) CreateAddress(userID uint, req *models.UserAddressReq) (uint, error) {
	address := models.UserAddress{
		UserID:     userID,
		UserName:   req.UserName,
		UserPhone:  req.UserPhone,
		Province:   req.Province,
		City:       req.City,
		District:   req.District,
		Detail:     req.Detail,
		PostalCode: req.PostalCode,
		IsDefault:  req.IsDefault,
	}

	// 如果设置为默认地址，取消其他默认地址
	if req.IsDefault == 1 {
		global.DB.Model(&models.UserAddress{}).Where("user_id = ?", userID).Update("is_default", 0)
	}

	if err := global.DB.Create(&address).Error; err != nil {
		return 0, err
	}
	return address.ID, nil
}

// UpdateAddress 更新地址
func (s *UserService) UpdateAddress(id, userID uint, req *models.UserAddressReq) error {
	// 如果设置为默认地址，取消其他默认地址
	if req.IsDefault == 1 {
		global.DB.Model(&models.UserAddress{}).Where("user_id = ?", userID).Update("is_default", 0)
	}

	return global.DB.Model(&models.UserAddress{}).Where("id = ? AND user_id = ?", id, userID).Updates(map[string]interface{}{
		"user_name":   req.UserName,
		"user_phone":  req.UserPhone,
		"province":    req.Province,
		"city":        req.City,
		"district":    req.District,
		"detail":      req.Detail,
		"postal_code": req.PostalCode,
		"is_default":  req.IsDefault,
	}).Error
}

// DeleteAddress 删除地址
func (s *UserService) DeleteAddress(id, userID uint) error {
	return global.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.UserAddress{}).Error
}

// ==================== 购物车管理 ====================

// GetCartList 获取购物车列表
func (s *UserService) GetCartList(userID uint) ([]models.UserCart, error) {
	var carts []models.UserCart
	err := global.DB.Where("user_id = ?", userID).Find(&carts).Error
	return carts, err
}

// GetCartByID 根据ID获取购物车项
func (s *UserService) GetCartByID(id, userID uint) (*models.UserCart, error) {
	var cart models.UserCart
	err := global.DB.Where("id = ? AND user_id = ?", id, userID).First(&cart).Error
	return &cart, err
}

// GetCartByProductAndSku 根据商品和SKU获取购物车项
func (s *UserService) GetCartByProductAndSku(userID uint, productID, skuID uint) (*models.UserCart, error) {
	var cart models.UserCart
	err := global.DB.Where("user_id = ? AND product_id = ? AND sku_id = ?", userID, productID, skuID).First(&cart).Error
	return &cart, err
}

// AddOrUpdateCart 添加或更新购物车
func (s *UserService) AddOrUpdateCart(userID uint, req *models.UserCartReq) (uint, error) {
	// 检查是否已存在
	var cart models.UserCart
	err := global.DB.Where("user_id = ? AND product_id = ? AND sku_id = ?", userID, req.ProductID, req.SkuID).First(&cart).Error

	if err == nil {
		// 更新数量
		cart.Quantity += req.Quantity
		if err := global.DB.Save(&cart).Error; err != nil {
			return 0, err
		}
		return cart.ID, nil
	}

	// 创建新记录
	cart = models.UserCart{
		UserID:    userID,
		ProductID: req.ProductID,
		SkuID:     req.SkuID,
		Quantity:  req.Quantity,
		Selected:  req.Selected,
	}
	if err := global.DB.Create(&cart).Error; err != nil {
		return 0, err
	}
	return cart.ID, nil
}

// UpdateCart 更新购物车
func (s *UserService) UpdateCart(id, userID uint, req *models.UserCartUpdateReq) error {
	return global.DB.Model(&models.UserCart{}).Where("id = ? AND user_id = ?", id, userID).Updates(map[string]interface{}{
		"quantity": req.Quantity,
		"selected": req.Selected,
	}).Error
}

// DeleteCart 删除购物车
func (s *UserService) DeleteCart(id, userID uint) error {
	return global.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.UserCart{}).Error
}
