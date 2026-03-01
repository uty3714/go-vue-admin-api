package flag

import (
	"flag"
	"fmt"
)

type Option struct {
	DB       bool
	ResetDB  bool
	ResetPwd bool
	Help     bool
}

// Parse 解析命令行参数
func Parse() Option {
	var opt Option
	flag.BoolVar(&opt.DB, "db", false, "初始化数据库")
	flag.BoolVar(&opt.ResetDB, "reset-db", false, "重置数据库（删除所有表并重新初始化）")
	flag.BoolVar(&opt.ResetPwd, "reset-pwd", false, "重置管理员密码为 admin123")
	flag.BoolVar(&opt.Help, "h", false, "帮助")
	flag.Parse()
	return opt
}

// Usage 打印帮助信息
func Usage() {
	fmt.Println("Go-Vue-Admin API 服务")
	fmt.Println()
	fmt.Println("用法:")
	fmt.Println("  ./admin-server [选项]")
	fmt.Println()
	fmt.Println("选项:")
	fmt.Println("  -db          初始化数据库（创建表和基础数据）")
	fmt.Println("  -reset-db    重置数据库（删除所有表并重新初始化，谨慎使用）")
	fmt.Println("  -reset-pwd   重置管理员密码为 admin123")
	fmt.Println("  -h           显示帮助信息")
	fmt.Println()
	fmt.Println("示例:")
	fmt.Println("  ./admin-server -db            # 初始化数据库")
	fmt.Println("  ./admin-server -reset-db      # 重置数据库")
	fmt.Println("  ./admin-server -reset-pwd     # 重置管理员密码")
	fmt.Println("  ./admin-server                # 启动服务器")
}
