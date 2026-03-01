package middleware

import (
	"strings"
	"go-vue-admin/models/res"
	"go-vue-admin/util"

	"github.com/gin-gonic/gin"
)

// MiniAuth 小程序认证中间件
func MiniAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取token
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			res.Unauthorized(c, "请求未携带token，无法访问")
			c.Abort()
			return
		}

		// 检查token格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			res.FailWithMessage(c, res.ErrorCodeTokenInvalid, "token格式错误")
			c.Abort()
			return
		}

		// 解析token
		j := util.NewJWT()
		claims, err := j.ParseToken(parts[1])
		if err != nil {
			if err == util.TokenExpired {
				res.FailWithMessage(c, res.ErrorCodeTokenExpired, "token已过期，请重新登录")
				c.Abort()
				return
			}
			res.FailWithMessage(c, res.ErrorCodeTokenInvalid, err.Error())
			c.Abort()
			return
		}

		// 设置用户信息
		c.Set("userId", claims.UserID)
		c.Set("openId", claims.Username) // 小程序用openId作为username
		c.Set("claims", claims)

		c.Next()
	}
}

// MiniAuthOptional 小程序可选认证（不强制要求登录）
func MiniAuthOptional() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取token
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		// 检查token格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.Next()
			return
		}

		// 解析token
		j := util.NewJWT()
		claims, err := j.ParseToken(parts[1])
		if err != nil {
			c.Next()
			return
		}

		// 设置用户信息
		c.Set("userId", claims.UserID)
		c.Set("openId", claims.Username)
		c.Set("claims", claims)

		c.Next()
	}
}
