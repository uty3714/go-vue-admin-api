package util

import (
	"go-vue-admin/global"

	"golang.org/x/crypto/bcrypt"
)

// BcryptHash 使用bcrypt对密码进行加密
func BcryptHash(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		global.Log.Errorf("密码哈希失败: %v", err)
		return ""
	}
	return string(bytes)
}

// BcryptCheck 对比明文密码和数据库的哈希值
func BcryptCheck(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
