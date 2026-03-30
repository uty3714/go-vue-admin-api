package middleware

import (
	"net/http"
	"net/url"
	"go-vue-admin/global"

	"github.com/gin-gonic/gin"
)

// Cors 跨域中间件（默认使用白名单模式，禁止反射任意Origin）
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		
		// 检查是否配置了白名单
		cfg := global.Config.Cors
		var isAllowed bool
		
		if cfg.AllowAll {
			// 如果明确配置了允许所有，使用通配符（但不允许携带凭证）
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Credentials", "false")
			isAllowed = true
		} else if len(cfg.Whitelist) > 0 {
			// 白名单模式
			for _, whitelist := range cfg.Whitelist {
				if isValidOrigin(origin, whitelist.AllowOrigin) != "" {
					c.Header("Access-Control-Allow-Origin", whitelist.AllowOrigin)
					if whitelist.AllowCredentials {
						c.Header("Access-Control-Allow-Credentials", "true")
					}
					isAllowed = true
					break
				}
			}
		}
		
		// 如果不在白名单中，且不是OPTIONS预检请求，拒绝访问
		if !isAllowed && cfg.Mode == "strict-whitelist" && method != "OPTIONS" {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		
		// 如果没有匹配到白名单，但不强制严格模式，则不允许携带凭证
		if !isAllowed {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Credentials", "false")
		}
		
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE, PATCH")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization, X-Token")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type, Authorization, X-Refresh-Token")
		c.Header("Access-Control-Max-Age", "86400")

		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// isValidOrigin 验证Origin是否匹配白名单
// 支持精确匹配和后缀匹配（如 https://*.example.com）
func isValidOrigin(origin, pattern string) string {
	if origin == "" {
		return ""
	}
	
	// 精确匹配
	if origin == pattern {
		return origin
	}
	
	// 解析origin
	originURL, err := url.Parse(origin)
	if err != nil {
		return ""
	}
	
	// 处理通配符模式（如 https://*.example.com）
	if len(pattern) > 2 && pattern[0] == '*' && pattern[1] == '.' {
		suffix := pattern[1:] // .example.com
		if len(originURL.Host) >= len(suffix) && 
		   originURL.Host[len(originURL.Host)-len(suffix):] == suffix {
			return origin
		}
	}
	
	return ""
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
			if isValidOrigin(origin, whitelist.AllowOrigin) != "" {
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
