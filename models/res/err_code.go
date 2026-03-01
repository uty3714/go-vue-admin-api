package res

// 响应状态码常量
const (
	// 成功
	SuccessCode = 200
	SuccessMsg  = "success"

	// 请求参数错误 1000-1999
	ErrorCodeParamInvalid = 1000
	ErrorMsgParamInvalid  = "请求参数错误"

	// 认证授权错误 2000-2999
	ErrorCodeUnauthorized     = 2000
	ErrorMsgUnauthorized      = "未登录或登录已过期"
	ErrorCodeForbidden        = 2001
	ErrorMsgForbidden         = "无权限访问"
	ErrorCodeTokenInvalid     = 2002
	ErrorMsgTokenInvalid      = "无效的token"
	ErrorCodeTokenExpired     = 2003
	ErrorMsgTokenExpired      = "token已过期"
	ErrorCodeLoginFailed      = 2004
	ErrorMsgLoginFailed       = "用户名或密码错误"
	ErrorCodeUserDisabled     = 2005
	ErrorMsgUserDisabled      = "用户已被禁用"
	ErrorCodeUserNotExist     = 2006
	ErrorMsgUserNotExist      = "用户不存在"
	ErrorCodeUserExist        = 2007
	ErrorMsgUserExist         = "用户已存在"
	ErrorCodePasswordError    = 2008
	ErrorMsgPasswordError     = "密码错误"
	ErrorCodeOldPasswordError = 2009
	ErrorMsgOldPasswordError  = "旧密码错误"

	// 业务逻辑错误 3000-3999
	ErrorCodeBusinessError = 3000
	ErrorMsgBusinessError  = "业务处理失败"

	// 资源不存在错误 4000-4999
	ErrorCodeNotFound = 4000
	ErrorMsgNotFound  = "资源不存在"

	// 服务器内部错误 5000-5999
	ErrorCodeInternalServer = 5000
	ErrorMsgInternalServer  = "服务器内部错误"
	ErrorCodeDBError        = 5001
	ErrorMsgDBError         = "数据库操作失败"
)

// ErrorCodeMap 错误码映射
var ErrorCodeMap = map[int]string{
	SuccessCode:               SuccessMsg,
	ErrorCodeParamInvalid:     ErrorMsgParamInvalid,
	ErrorCodeUnauthorized:     ErrorMsgUnauthorized,
	ErrorCodeForbidden:        ErrorMsgForbidden,
	ErrorCodeTokenInvalid:     ErrorMsgTokenInvalid,
	ErrorCodeTokenExpired:     ErrorMsgTokenExpired,
	ErrorCodeLoginFailed:      ErrorMsgLoginFailed,
	ErrorCodeUserDisabled:     ErrorMsgUserDisabled,
	ErrorCodeUserNotExist:     ErrorMsgUserNotExist,
	ErrorCodeUserExist:        ErrorMsgUserExist,
	ErrorCodePasswordError:    ErrorMsgPasswordError,
	ErrorCodeOldPasswordError: ErrorMsgOldPasswordError,
	ErrorCodeBusinessError:    ErrorMsgBusinessError,
	ErrorCodeNotFound:         ErrorMsgNotFound,
	ErrorCodeInternalServer:   ErrorMsgInternalServer,
	ErrorCodeDBError:          ErrorMsgDBError,
}

// GetErrorMsg 获取错误信息
func GetErrorMsg(code int) string {
	if msg, ok := ErrorCodeMap[code]; ok {
		return msg
	}
	return "未知错误"
}
