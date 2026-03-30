package middleware

import (
	"go-vue-admin/global"
	"go-vue-admin/models/res"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorHandler 统一错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			global.Log.Errorf("请求处理错误: %v", err)

			// 根据错误类型返回不同的响应
			switch err.Type {
			case gin.ErrorTypeBind:
				res.ValidationError(c, "参数绑定失败: "+err.Err.Error())
			case gin.ErrorTypeRender:
				res.FailWithMessage(c, res.ErrorCodeInternalServer, "渲染失败")
			default:
				res.Error(c, err.Err)
			}
			return
		}

		// 处理HTTP状态码错误
		if c.Writer.Status() != http.StatusOK && c.Writer.Status() != http.StatusNoContent {
			switch c.Writer.Status() {
			case http.StatusNotFound:
				res.Fail(c, res.ErrorCodeNotFound)
			case http.StatusForbidden:
				res.Fail(c, res.ErrorCodeForbidden)
			case http.StatusUnauthorized:
				res.Fail(c, res.ErrorCodeUnauthorized)
			case http.StatusInternalServerError:
				res.Fail(c, res.ErrorCodeInternalServer)
			}
		}
	}
}
