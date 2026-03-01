package middleware

import (
	"strings"
	"go-vue-admin/global"
	"go-vue-admin/models/res"
	"go-vue-admin/util"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT认证中间件
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

		// 检查用户状态
		var userId = claims.UserID
		c.Set("userId", userId)
		c.Set("username", claims.Username)
		c.Set("roleId", claims.RoleID)
		c.Set("claims", claims)

		global.Log.Infof("用户[%s]访问: %s %s", claims.Username, c.Request.Method, c.Request.URL.Path)

		c.Next()
	}
}

// JWTAuthWithRole 带角色验证的JWT认证中间件
func JWTAuthWithRole(roleCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		JWTAuth()(c)
		
		if c.IsAborted() {
			return
		}

		// TODO: 根据角色代码验证权限
		roleId, _ := c.Get("roleId")
		global.Log.Infof("验证角色权限, roleId: %v, required: %s", roleId, roleCode)

		c.Next()
	}
}
