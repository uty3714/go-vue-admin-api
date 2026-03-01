package core

import (
	"os"
	"path"
	"time"
	"go-vue-admin/global"

	"github.com/sirupsen/logrus"
)

func InitLogrus() *logrus.Logger {
	m := global.Config.Zap

	logger := logrus.New()

	// 设置日志级别
	level, err := logrus.ParseLevel(m.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	// 设置日志格式
	if m.Format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
			ForceColors:     true,
		})
	}

	// 创建日志目录
	if _, err := os.Stat(m.Director); os.IsNotExist(err) {
		os.MkdirAll(m.Director, 0755)
	}

	// 打开日志文件
	logPath := path.Join(m.Director, time.Now().Format("2006-01-02")+".log")
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.Error("打开日志文件失败: ", err)
	}

	// 设置输出
	if m.LogInConsole {
		logger.SetOutput(&dualWriter{
			file:   file,
			stdout: os.Stdout,
		})
	} else {
		logger.SetOutput(file)
	}

	global.Log = logger
	return logger
}

type dualWriter struct {
	file   *os.File
	stdout *os.File
}

func (w *dualWriter) Write(p []byte) (n int, err error) {
	w.file.Write(p)
	return w.stdout.Write(p)
}
