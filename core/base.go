package core

import (
	"context"
	"encoding/json"
	"sync"
	"time"
	"zpic-client/helper"

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

var appCtx context.Context

// 设置全局上下文
func SetCtx(ctx context.Context) {
	appCtx = ctx
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

// UploadReq 上传请求参数
type UploadReq struct {
	FilePath string `json:"file_path"`
	AlbumID  int64  `json:"album_id"`
}

// UploadResData 上传成功后返回的图片信息
type UploadResData struct {
	ImgID        string `json:"imgid"`
	Path         string `json:"path"`
	URL          string `json:"url"`
	ThumbnailURL string `json:"thumbnail_url"`
	Width        int64  `json:"width"`
	Height       int64  `json:"height"`
	Filename     string `json:"filename"`
	Size         int64  `json:"size"`
}

// uploadParams 上传时附带的额外参数
type uploadParams struct {
	Dedup     bool  `json:"dedup"`
	AlbumID   int64 `json:"album_id"`
	Watermark bool  `json:"watermark"`
	Compress  bool  `json:"compress"`
}

var (
	uploadClient *resty.Client
	uploadOnce   sync.Once
)

// getUploadClient 返回上传专用的 resty.Client 单例，避免多次调用重复创建。
func getUploadClient() *resty.Client {
	uploadOnce.Do(func() {
		baseURL := viper.GetString("base_url")
		token := viper.GetString("token")
		uploadClient = resty.New().
			SetHeader("User-Agent", "zpic-client/1.0.0").
			SetTimeout(120*time.Second).
			SetBaseURL(baseURL+"/api/v3/upload").
			SetHeader("Authorization", "Bearer "+token)
	})
	return uploadClient
}

// UploadZpic 上传图片到 Zpic，返回是否成功以及上传后的图片信息。
func UploadZpic(req UploadReq) (bool, UploadResData) {
	client := getUploadClient()

	// 构造 params JSON 字符串
	params := uploadParams{
		Dedup:     true,
		AlbumID:   req.AlbumID,
		Watermark: false,
		Compress:  false,
	}
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		helper.WriteLog("序列化 params 失败：" + err.Error())
		return false, UploadResData{}
	}

	resp, err := client.R().
		SetFile("file", req.FilePath).
		SetFormData(map[string]string{
			"params": string(paramsJSON),
		}).
		Post("")
	if err != nil {
		helper.WriteLog("上传请求失败：" + err.Error())
		return false, UploadResData{}
	}

	bodyBytes := resp.Body()
	code := gjson.GetBytes(bodyBytes, "code").Int()
	if code != 200 {
		msg := gjson.GetBytes(bodyBytes, "msg").String()
		helper.WriteLog("上传失败：" + msg)
		return false, UploadResData{}
	}

	data := gjson.GetBytes(bodyBytes, "data")
	result := UploadResData{
		ImgID:        data.Get("imgid").String(),
		Path:         data.Get("path").String(),
		URL:          data.Get("url").String(),
		ThumbnailURL: data.Get("thumbnail_url").String(),
		Width:        data.Get("width").Int(),
		Height:       data.Get("height").Int(),
		Filename:     data.Get("filename").String(),
		Size:         data.Get("size").Int(),
	}

	return true, result
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
