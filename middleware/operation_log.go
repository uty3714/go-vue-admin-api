package middleware

import (
	"bytes"
	"go-vue-admin/global"
	"go-vue-admin/models"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// 操作日志工作池，限制并发goroutine数量防止泄露
var (
	logChan   = make(chan models.OperationLog, 1000) // 缓冲通道
	logOnce   sync.Once
)

// initLogWorker 初始化日志工作池
func initLogWorker() {
	logOnce.Do(func() {
		// 启动固定数量的工作协程
		for i := 0; i < 5; i++ {
			go func() {
				for log := range logChan {
					if err := global.DB.Create(&log).Error; err != nil {
						global.Log.Errorf("保存操作日志失败: %v", err)
					}
				}
			}()
		}
	})
}

func init() {
	initLogWorker()
}

// OperationLog 操作日志中间件
func OperationLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 获取请求数据
		var requestData string
		
		// GET 请求记录查询参数，其他请求记录 Body
		if c.Request.Method == "GET" {
			query := c.Request.URL.Query().Encode()
			if query != "" {
				requestData = "[Query] " + query
			}
		} else if c.Request.Body != nil {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err != nil {
				global.Log.Warnf("读取请求体失败: %v", err)
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			// 过滤敏感字段
			requestData = filterSensitiveData(string(bodyBytes))
		}

		// 创建自定义ResponseWriter来捕获响应
		blw := &bodyLogWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
			mu:             &sync.Mutex{},
		}
		c.Writer = blw

		// 处理请求
		c.Next()

		// 计算耗时
		duration := time.Since(startTime).Milliseconds()

		// 获取用户信息（安全类型断言）
		userId, _ := c.Get("userId")
		username, _ := c.Get("username")
		roleId, _ := c.Get("roleId")

		// 获取角色名称（安全类型断言）
		roleName := ""
		if roleId != nil {
			if rid, ok := roleId.(uint); ok {
				var role models.SystemRole
				if err := global.DB.First(&role, rid).Error; err != nil {
					global.Log.Warnf("查询角色失败: %v", err)
				} else {
					roleName = role.RoleName
				}
			}
		}

		// 确定操作状态
		status := 1 // 成功
		if c.Writer.Status() >= 400 {
			status = 2 // 失败
		}

		// 过滤响应数据
		responseData := blw.body.String()
		
		// 对日志查询接口本身，不记录响应数据（避免套娃）
		path := c.Request.URL.Path
		if path == "/api/v1/system/log/operation/list" || path == "/api/v1/system/log/login/list" {
			responseData = "[日志列表数据省略]"
		} else if len(responseData) > 1000 {
			// 其他接口限制响应数据长度
			responseData = responseData[:1000] + "..."
		}

		// 保存操作日志
		log := models.OperationLog{
			UserID:        getUint(userId),
			Username:      getString(username),
			RoleName:      roleName,
			Method:        c.Request.Method,
			Path:          path,
			RequestData:   requestData,
			ResponseData:  responseData,
			Status:        status,
			ErrorMessage:  getErrorMessage(c),
			IP:            c.ClientIP(),
			UserAgent:     c.Request.UserAgent(),
			OperationTime: int(duration),
			CreatedAt:     models.LocalTime(time.Now()),
		}

		// 发送到日志通道（有界队列，防止goroutine泄露）
		select {
		case logChan <- log:
			// 成功发送到队列
		default:
			// 队列已满，丢弃日志并记录警告
			global.Log.Warn("操作日志队列已满，丢弃日志记录")
		}
	}
}

// bodyLogWriter 自定义ResponseWriter（线程安全）
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
	mu   *sync.Mutex
}

// Write 实现Write方法（线程安全）
func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.mu.Lock()
	w.body.Write(b)
	w.mu.Unlock()
	return w.ResponseWriter.Write(b)
}

// filterSensitiveData 过滤敏感数据
func filterSensitiveData(data string) string {
	// 过滤密码字段
	sensitiveFields := []string{"password", "oldPassword", "newPassword", "confirmPassword"}
	for _, field := range sensitiveFields {
		// 简单替换，实际生产环境可以使用更复杂的正则
		if strings.Contains(data, field) {
			data = "[FILTERED]"
			break
		}
	}
	return data
}

// getUint 安全获取uint
func getUint(v interface{}) uint {
	if v == nil {
		return 0
	}
	if id, ok := v.(uint); ok {
		return id
	}
	return 0
}

// getString 安全获取string
func getString(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

// getErrorMessage 获取错误信息
func getErrorMessage(c *gin.Context) string {
	if len(c.Errors) > 0 {
		return c.Errors.String()
	}
	return ""
}
