package core

import (
	"github.com/spf13/viper"
)

// 声明一个结构体，用于存储前端传递过来的设置参数
type SettingData struct {
	// 基础域名
	BaseUrl string `json:"base_url"`
	// API token
	Token string `json:"token"`
	// http代理
	HttpProxy string `json:"http_proxy"`
}

func (ac *AppCore) GetSetting() SettingData {
	return SettingData{
		BaseUrl:   viper.GetString("base_url"),
		Token:     viper.GetString("token"),
		HttpProxy: viper.GetString("http_proxy"),
	}
}

func (ac *AppCore) UpdateSetting(data SettingData) bool {
	// 遍历结构体字段，更新配置
	viper.Set("base_url", data.BaseUrl)
	viper.Set("token", data.Token)
	viper.Set("http_proxy", data.HttpProxy)
	// 写入配置并保存
	err := viper.WriteConfig()
	if err != nil {
		return false
	}
	return true
}
