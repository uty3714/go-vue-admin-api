package middleware

import (
	"strings"
	"time"
	"go-vue-admin/global"
	"go-vue-admin/models"
	"go-vue-admin/models/res"
	"go-vue-admin/util"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT认证中间件
// 支持Token自动刷新，刷新后的旧Token会被加入黑名单
func JWTAuth() gin.HandlerFunc {
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

		tokenString := parts[1]
		j := util.NewJWT()
		claims, err := j.ParseToken(tokenString)
		
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

		// 检查是否需要刷新token（在过期前1小时内）
		if claims.ExpiresAt != nil {
			bufferTime := time.Duration(global.Config.JWT.BufferTime) * time.Hour
			if time.Until(claims.ExpiresAt.Time) < bufferTime {
				// 刷新token
				newToken, err := j.RefreshToken(tokenString)
				if err == nil {
					c.Header("X-Refresh-Token", newToken)
					global.Log.Infof("用户[%s]的token已自动刷新", claims.Username)
					
					// 将旧token加入黑名单，防止重用攻击
					tb := &TokenBlacklist{}
					if err := tb.AddToBlacklist(tokenString, claims.ExpiresAt.Time); err != nil {
						global.Log.Errorf("将旧token加入黑名单失败: %v", err)
					}
				}
			}
		}

		// 设置上下文值
		c.Set("userId", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("roleId", claims.RoleID)
		c.Set("claims", claims)

		global.Log.Infof("用户[%s]访问: %s %s", claims.Username, c.Request.Method, c.Request.URL.Path)

		c.Next()
	}
}

// JWTAuthWithRole 带角色验证的JWT认证中间件
// 需要指定角色代码才能访问
func JWTAuthWithRole(requiredRoleCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		JWTAuth()(c)
		
		if c.IsAborted() {
			return
		}

		roleId, exists := c.Get("roleId")
		if !exists {
			res.Fail(c, res.ErrorCodeForbidden)
			c.Abort()
			return
		}

		// 安全的类型断言
		roleIdUint, ok := roleId.(uint)
		if !ok {
			global.Log.Errorf("roleId 类型断言失败: %T", roleId)
			res.Fail(c, res.ErrorCodeForbidden)
			c.Abort()
			return
		}

		// 查询用户角色代码
		var role models.SystemRole
		if err := global.DB.First(&role, roleIdUint).Error; err != nil {
			global.Log.Errorf("查询角色失败: %v", err)
			res.Fail(c, res.ErrorCodeForbidden)
			c.Abort()
			return
		}

		// 检查是否是超级管理员（admin有所有权限）
		if role.RoleCode == "admin" {
			c.Next()
			return
		}

		// 验证角色代码是否匹配
		if role.RoleCode != requiredRoleCode {
			global.Log.Warnf("角色权限不足, userRole: %s, required: %s", role.RoleCode, requiredRoleCode)
			res.Fail(c, res.ErrorCodeForbidden)
			c.Abort()
			return
		}

		c.Next()
	}
}
