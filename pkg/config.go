package pkg

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"zpic-client/helper"

	"github.com/spf13/viper"
)

// 初始化配置文件
func InitConfig() {
	runDir := helper.GetRunDir()
	// 获取系统临时目录
	tempDir := os.TempDir()
	tempDir = filepath.Join(tempDir, "zpic-temp")
	// 创建必要的目录
	dirs := []string{runDir + "/data", runDir + "/data/config", runDir + "/data/db", runDir + "/data/logs", runDir + "/data/temp", tempDir}
	for _, dir := range dirs {
		err := helper.CreateDir(dir)
		if err != nil {
			// 打印错误
			fmt.Printf("%s\n", err)
			// 退出执行
			os.Exit(1)
		}
	}

	//配置文件目录
	config_dir := runDir + "/data/config"
	//配置文件路径
	config_file := config_dir + "/config.toml"
	//检查配置文件是否存在，如果存在了，则不进行初始化
	_, err := os.Stat(config_file)
	//终止更新脚本
	// KillZdirUpdater()
	//返回的error为空，说明文件存在，存在则不允许再次初始化,直接就返回了
	if err == nil {
		return
	} else {
		//如果配置文件目录不存在，则创建
		if !helper.V_dir(config_dir) {
			err := helper.CreateDir(config_dir)

			if err != nil {
				fmt.Printf("%s\n", err)
				os.Exit(1)
			}
		}

		//创建目标文件
		target, t_error := os.Create(config_file)
		if t_error != nil {
			fmt.Printf("%s\n", t_error)
			os.Exit(1)
		}

		defer target.Close()

	}
}

var once sync.Once

// 加载配置文件
func LoadConfig() {
	// 初始话配置文件
	InitConfig()
	runDir := helper.GetRunDir()
	once.Do(func() {
		//默认配置文件
		config_file := runDir + "/data/config/config.toml"

		viper.SetConfigFile(config_file) // 指定配置文件路径
		//指定ini类型的文件
		viper.SetConfigType("toml")
		err := viper.ReadInConfig() // 读取配置信息
		if err != nil {             // 读取配置信息失败
			// 写入日志
			helper.WriteLog(fmt.Sprintf("%s", err))
			fmt.Println("Failed to read config:", err)
		}
	})

}
