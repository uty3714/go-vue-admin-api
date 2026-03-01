package util

import (
	"os"
	"path/filepath"
)

// PathExists 判断文件或目录是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// CreateDir 创建目录
func CreateDir(path string) error {
	if exist, _ := PathExists(path); !exist {
		return os.MkdirAll(path, os.ModePerm)
	}
	return nil
}

// GetCurrentDirectory 获取当前目录
func GetCurrentDirectory() string {
	dir, _ := os.Getwd()
	return dir
}

// GetFileSize 获取文件大小
func GetFileSize(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return info.Size()
}

// GetFileExt 获取文件扩展名
func GetFileExt(filename string) string {
	return filepath.Ext(filename)
}

// GetFileName 获取文件名（不含扩展名）
func GetFileName(filename string) string {
	ext := filepath.Ext(filename)
	return filename[:len(filename)-len(ext)]
}

// GetFileNameWithExt 获取文件名（含扩展名）
func GetFileNameWithExt(path string) string {
	return filepath.Base(path)
}
