package util

import (
	"strconv"
)

// StringToInt string转int
func StringToInt(str string) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return num
}

// StringToUint string转uint
func StringToUint(str string) uint {
	num, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0
	}
	return uint(num)
}

// IntToString int转string
func IntToString(num int) string {
	return strconv.Itoa(num)
}

// UintToString uint转string
func UintToString(num uint) string {
	return strconv.FormatUint(uint64(num), 10)
}

// StringToInt64 string转int64
func StringToInt64(str string) int64 {
	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	}
	return num
}

// Int64ToString int64转string
func Int64ToString(num int64) string {
	return strconv.FormatInt(num, 10)
}
