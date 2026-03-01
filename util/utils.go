package util

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"time"
)

// MD5 md5加密
func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// MD5V md5加密（盐值）
func MD5V(str string, salt string) string {
	return MD5(str + salt)
}

// GenerateUUID 生成UUID
func GenerateUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// GenerateRandomString 生成随机字符串
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[num.Int64()]
	}
	return string(result)
}

// GenerateRandomNumber 生成随机数字
func GenerateRandomNumber(length int) string {
	const charset = "0123456789"
	result := make([]byte, length)
	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[num.Int64()]
	}
	return string(result)
}

// GenerateOrderNo 生成订单号
func GenerateOrderNo() string {
	return time.Now().Format("20060102150405") + GenerateRandomNumber(6)
}

// IsEmail 验证邮箱格式
func IsEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

// IsMobile 验证手机号格式（中国大陆）
func IsMobile(mobile string) bool {
	pattern := `^1[3-9]\d{9}$`
	matched, _ := regexp.MatchString(pattern, mobile)
	return matched
}

// IsPhone 验证电话号码格式
func IsPhone(phone string) bool {
	pattern := `^\d{3,4}-?\d{7,8}$`
	matched, _ := regexp.MatchString(pattern, phone)
	return matched
}

// Trim 去除字符串首尾空格
func Trim(str string) string {
	return strings.TrimSpace(str)
}

// Contains 判断字符串是否包含子串
func Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// ToLower 字符串转小写
func ToLower(str string) string {
	return strings.ToLower(str)
}

// ToUpper 字符串转大写
func ToUpper(str string) string {
	return strings.ToUpper(str)
}

// SubString 截取字符串
func SubString(str string, start, length int) string {
	runes := []rune(str)
	if start < 0 {
		start = 0
	}
	if start >= len(runes) {
		return ""
	}
	end := start + length
	if end > len(runes) {
		end = len(runes)
	}
	return string(runes[start:end])
}

// FormatTime 格式化时间
func FormatTime(t time.Time, format string) string {
	return t.Format(format)
}

// ParseTime 解析时间字符串
func ParseTime(timeStr string, format string) (time.Time, error) {
	return time.Parse(format, timeStr)
}

// GetCurrentTime 获取当前时间
func GetCurrentTime() time.Time {
	return time.Now()
}

// GetTimestamp 获取当前时间戳（秒）
func GetTimestamp() int64 {
	return time.Now().Unix()
}

// GetTimestampMs 获取当前时间戳（毫秒）
func GetTimestampMs() int64 {
	return time.Now().UnixNano() / 1e6
}
