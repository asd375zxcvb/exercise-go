package main

import (
	"fmt"
	"github.com/spf13/viper"
)

func main() {
	// 初始化viper
	viper.SetConfigName("config") // 指定配置文件名（不带后缀）
	viper.AddConfigPath(".")      // 搜索路径
	err := viper.ReadInConfig()   // 查找并读取配置文件
	if err != nil {               // 处理读取错误
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// 从配置文件读取配置项
	host := viper.GetString("database.host")
	port := viper.GetInt("database.port")
	username := viper.GetString("database.username")
	password := viper.GetString("database.password")

	// 打印配置项
	fmt.Printf("host=%s, port=%d, username=%s, password=%s\n", host, port, username, password)
}
