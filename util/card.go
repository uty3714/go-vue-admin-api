package util

import (
	"regexp"
	"strconv"
	"time"
)

// ValidateIDCard 验证身份证号格式
func ValidateIDCard(idCard string) bool {
	if len(idCard) != 18 {
		return false
	}
	// 简单正则验证
	pattern := `^\d{17}[\dXx]$`
	matched, _ := regexp.MatchString(pattern, idCard)
	if !matched {
		return false
	}
	return validateIDCardChecksum(idCard)
}

// validateIDCardChecksum 验证身份证校验码
func validateIDCardChecksum(idCard string) bool {
	weight := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	checkCode := []string{"1", "0", "X", "9", "8", "7", "6", "5", "4", "3", "2"}

	sum := 0
	for i := 0; i < 17; i++ {
		num, _ := strconv.Atoi(string(idCard[i]))
		sum += num * weight[i]
	}

	mod := sum % 11
	return string(idCard[17]) == checkCode[mod] || (idCard[17] == 'x' && checkCode[mod] == "X")
}

// GetBirthFromIDCard 从身份证号获取出生日期
func GetBirthFromIDCard(idCard string) string {
	if len(idCard) != 18 {
		return ""
	}
	return idCard[6:10] + "-" + idCard[10:12] + "-" + idCard[12:14]
}

// GetAgeFromIDCard 从身份证号获取年龄
func GetAgeFromIDCard(idCard string) int {
	if len(idCard) != 18 {
		return 0
	}
	birthYear, _ := strconv.Atoi(idCard[6:10])
	currentYear := time.Now().Year()
	return currentYear - birthYear
}

// GetGenderFromIDCard 从身份证号获取性别
func GetGenderFromIDCard(idCard string) string {
	if len(idCard) != 18 {
		return ""
	}
	genderCode, _ := strconv.Atoi(idCard[16:17])
	if genderCode%2 == 0 {
		return "女"
	}
	return "男"
}
