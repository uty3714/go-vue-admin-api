package middleware

import (
	"net/http"
	"go-vue-admin/global"

	"github.com/gin-gonic/gin"
)

// Cors 跨域中间件
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE, PATCH")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization, X-Token")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// CorsByRules 根据配置处理跨域
func CorsByRules() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := global.Config.Cors

		// 允许所有跨域
		if cfg.AllowAll {
			Cors()(c)
			return
		}

		origin := c.Request.Header.Get("Origin")
		method := c.Request.Method

		// 白名单检查
		var isAllowed bool
		for _, whitelist := range cfg.Whitelist {
			if whitelist.AllowOrigin == origin {
				c.Header("Access-Control-Allow-Origin", whitelist.AllowOrigin)
				if whitelist.AllowMethods != "" {
					c.Header("Access-Control-Allow-Methods", whitelist.AllowMethods)
				} else {
					c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
				}
				if whitelist.AllowHeaders != "" {
					c.Header("Access-Control-Allow-Headers", whitelist.AllowHeaders)
				} else {
					c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
				}
				if whitelist.ExposeHeaders != "" {
					c.Header("Access-Control-Expose-Headers", whitelist.ExposeHeaders)
				}
				if whitelist.AllowCredentials {
					c.Header("Access-Control-Allow-Credentials", "true")
				}
				isAllowed = true
				break
			}
		}

		if !isAllowed && cfg.Mode == "strict-whitelist" {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
