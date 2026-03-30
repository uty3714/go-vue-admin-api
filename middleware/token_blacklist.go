package middleware

import (
	"go-vue-admin/global"
	"go-vue-admin/models/res"
	"go-vue-admin/util"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// 确保表只创建一次
var initBlacklistOnce sync.Once

// TokenBlacklist Token黑名单管理
type TokenBlacklist struct{}

// initBlacklistTable 初始化黑名单表（只执行一次）
func initBlacklistTable() error {
	var initErr error
	initBlacklistOnce.Do(func() {
		sql := `
		CREATE TABLE IF NOT EXISTS token_blacklist (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			token VARCHAR(512) NOT NULL UNIQUE COMMENT 'Token字符串',
			expires_at DATETIME NOT NULL COMMENT 'Token过期时间',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '加入黑名单时间',
			INDEX idx_expires (expires_at)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Token黑名单表';
		`
		if err := global.DB.Exec(sql).Error; err != nil {
			initErr = err
			global.Log.Errorf("创建token黑名单表失败: %v", err)
		}
	})
	return initErr
}

// AddToBlacklist 将token加入黑名单
func (tb *TokenBlacklist) AddToBlacklist(token string, expiresAt time.Time) error {
	// 初始化表（只执行一次）
	if err := initBlacklistTable(); err != nil {
		return err
	}

	// 插入黑名单
	result := global.DB.Exec(
		"INSERT IGNORE INTO token_blacklist (token, expires_at) VALUES (?, ?)",
		token, expiresAt,
	)
	if result.Error != nil {
		global.Log.Errorf("添加token到黑名单失败: %v", result.Error)
		return result.Error
	}
	return nil
}

// IsBlacklisted 检查token是否在黑名单中
func (tb *TokenBlacklist) IsBlacklisted(token string) bool {
	var count int64
	result := global.DB.Raw(
		"SELECT COUNT(*) FROM token_blacklist WHERE token = ? AND expires_at > NOW()",
		token,
	).Scan(&count)
	
	// 如果查询出错（如表不存在），允许访问（表不存在意味着没有token被拉黑）
	if result.Error != nil {
		global.Log.Debugf("检查token黑名单时出错（表可能不存在）: %v", result.Error)
		return false
	}
	
	return count > 0
}

// CleanupExpired 清理过期的黑名单记录
// 建议定期调用（如每天一次）
func (tb *TokenBlacklist) CleanupExpired() error {
	return global.DB.Exec("DELETE FROM token_blacklist WHERE expires_at <= NOW()").Error
}

// TokenBlacklistMiddleware Token黑名单检查中间件
func TokenBlacklistMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		token := parts[1]
		tb := &TokenBlacklist{}

		if tb.IsBlacklisted(token) {
			res.Fail(c, res.ErrorCodeTokenInvalid)
			c.Abort()
			return
		}

		c.Next()
	}
}

// LogoutHandler 登出处理
func LogoutHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		res.Success(c, nil)
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		res.Success(c, nil)
		return
	}

	token := parts[1]

	// 解析token获取过期时间
	j := util.NewJWT()
	claims, err := j.ParseToken(token)
	if err != nil {
		// Token已过期或无效，直接返回成功
		res.Success(c, nil)
		return
	}

	// 将token加入黑名单
	tb := &TokenBlacklist{}
	if claims.ExpiresAt != nil {
		if err := tb.AddToBlacklist(token, claims.ExpiresAt.Time); err != nil {
			global.Log.Errorf("添加token到黑名单失败: %v", err)
			// 继续返回成功，因为登出操作本身已完成
		}
	}

	res.Success(c, nil)
}
