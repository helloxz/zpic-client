package core

import (
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

// 声明返回的通用结构体
type ResData struct {
	Status bool        `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}

// 封装resty/v2，自动携带Token和设置超时时间以及基础URL
type HttpParams struct {
	Path   string
	Method string
}

// 返回resty client实例
// ReqZpic 返回已配置好的 resty.Client（包含 base_url、超时与 Authorization 头）。
// 示例：
//
//	client := ReqZpic("/api/v1/ping")
//	resp, err := client.R().Get("")
//	// 或：client.R().SetBody(payload).Post("")
func ReqZpic(path string) *resty.Client {
	timeout := 60 * time.Second
	base_url := viper.GetString("base_url")
	token := viper.GetString("token")
	api_url := base_url + path
	// 设置UA
	client := resty.New().
		SetHeader("User-Agent", "zpic-client/1.0.0").
		SetTimeout(timeout).
		SetBaseURL(api_url).
		SetHeader("Authorization", "Bearer "+token)

	return client
}

// AlbumItem 相册条目
type AlbumItem struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// 获取相册列表
func (ac *AppCore) GetAlbumList() ResData {
	path := "/api/v3/album_list"

	client := ReqZpic(path)
	resp, err := client.R().Get("")
	if err != nil {
		return ResData{
			Status: false,
			Msg:    "请求失败：" + err.Error(),
			Data:   nil,
		}
	}

	bodyBytes := resp.Body()
	code := gjson.GetBytes(bodyBytes, "code").Int()
	msg := gjson.GetBytes(bodyBytes, "msg").String()
	itemsResult := gjson.GetBytes(bodyBytes, "data.items")
	if !itemsResult.Exists() {
		return ResData{
			Status: false,
			Msg:    "响应中未找到 data.items",
			Data:   nil,
		}
	}
	items := itemsResult.Array()
	status := true

	if code != 200 {
		status = false
		return ResData{
			Status: status,
			Msg:    msg,
			Data:   nil,
		}
	}

	albumList := make([]AlbumItem, 0, len(items))
	for _, item := range items {
		albumList = append(albumList, AlbumItem{
			ID:   item.Get("album_id").Int(),
			Name: item.Get("name").String(),
		})
	}

	return ResData{
		Status: status,
		Msg:    msg,
		Data:   albumList,
	}
}
