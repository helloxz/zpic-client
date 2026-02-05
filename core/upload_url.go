package core

import (
	"bufio"
	"fmt"
	"net/url"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
	"zpic-client/helper"
	"zpic-client/model"

	"os"

	"github.com/go-resty/resty/v2"
	"github.com/sourcegraph/conc/pool"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

var wg sync.WaitGroup

// 声明ids数据结构体
type IdsStatus struct {
	Ids    []uint `json:"ids"`
	Status int8   `json:"status"`
}

// 声明ids结构体
type IdsForm struct {
	Ids []uint `json:"ids"`
}

// 声明前端结构体
type UrlsForm struct {
	AlbumID int    `json:"album_id"`
	Urls    string `json:"urls"`
}

// 声明列表的结构体
type UrlsList struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// 获取URL列表，根据ID降序排序，支持分页
func (ac *AppCore) GetUrlsList(data UrlsList) ResData {
	page := data.Page
	// 如果page小于1，则设置为1
	if page < 1 {
		page = 1
	}
	limit := data.Limit
	// 计算offset
	offset := (page - 1) * limit

	var urls []model.ZPurls
	result := model.DB.Order("id desc").Limit(limit).Offset(offset).Find(&urls)
	if result.Error != nil {
		return ResData{
			Status: false,
			Msg:    "获取URL列表失败：" + result.Error.Error(),
			Data:   nil,
		}
	}

	// 获取总数
	var total int64
	model.DB.Model(&model.ZPurls{}).Count(&total)

	return ResData{
		Status: true,
		Msg:    "获取URL列表成功",
		Data: map[string]interface{}{
			"items": urls,
			"total": total,
		},
	}
}

// 批量更改id状态
func (ac *AppCore) UpdateUrlsStatus(data IdsStatus) ResData {
	ids := data.Ids
	status := data.Status

	// 批量更新
	result := model.DB.Model(&model.ZPurls{}).Where("id IN ?", ids).Update("status", status)
	if result.Error != nil {
		return ResData{
			Status: false,
			Msg:    "更新URL状态失败：" + result.Error.Error(),
			Data:   nil,
		}
	}

	return ResData{
		Status: true,
		Msg:    "更新URL状态成功",
		Data:   nil,
	}
}

// 获取URL列表总数
func (ac *AppCore) GetUrlsCount() int64 {
	var count int64
	model.DB.Model(&model.ZPurls{}).Count(&count)
	return count
}

// 添加URL任务
func (ac *AppCore) AddUrls(data UrlsForm) ResData {
	albumID := data.AlbumID
	// 相册ID不能小于0
	if albumID < 0 {
		return ResData{
			Status: false,
			Msg:    "相册ID错误",
			Data:   nil,
		}
	}
	// 将urls按行分割
	scanner := bufio.NewScanner(strings.NewReader(data.Urls))
	var lines []string
	// 逐行扫描URL列表
	for scanner.Scan() {
		text := scanner.Text()
		// 移除前后空格
		text = strings.TrimSpace(text)
		if text == "" {
			continue
		}
		// 判断是否是有效的URL格式，不是有效的URL则跳过
		if !IsValidURL(text) {
			continue
		}
		lines = append(lines, text)
	}
	// 获取长度
	urls_count := len(lines)
	// 不能超过100个
	if urls_count > 100 {
		return ResData{
			Status: false,
			Msg:    "一次最多添加100个URL，请拆分后再添加",
			Data:   nil,
		}
	}
	// 批量插入
	var urlsData []model.ZPurls
	for _, line := range lines {
		urlsData = append(urlsData, model.ZPurls{
			AlbumID:   albumID,
			OriginURL: line,
		})
	}
	// 插入数据库
	result := model.DB.Create(&urlsData)

	if result.Error != nil {
		return ResData{
			Status: false,
			Msg:    "添加URL任务失败：" + result.Error.Error(),
			Data:   nil,
		}
	}

	return ResData{
		Status: true,
		Msg:    fmt.Sprintf("成功添加%d条数据", urls_count),
		Data:   nil,
	}
}

type taskData struct {
	ID        uint
	AlbumID   int
	OriginURL string
}

// 任务列表
func UploadTaskList() {
	// 根据ID增序排序，查询出status为0的任务，限制10条，然后并发上传
	limit := 30
	var tasks []taskData
	result := model.DB.Model(&model.ZPurls{}).Select("id, album_id, origin_url").Where("status = ?", model.UploadStatusPending).Order("id asc").Limit(limit).Find(&tasks)
	if result.Error != nil {
		return
	}
	// 批量将查询到的数据状态改为上传中
	var ids []uint
	for _, task := range tasks {
		ids = append(ids, task.ID)
	}
	model.DB.Model(&model.ZPurls{}).Where("id IN ?", ids).Update("status", model.UploadStatusUploading)
	// 遍历任务列表，启动协程上传
	p := pool.New().WithMaxGoroutines(10)
	for _, task := range tasks {
		t := task
		p.Go(func() {
			UploadURL(t)
			// 暂停0.5s，避免过快
			time.Sleep(500 * time.Millisecond)
		})
	}
	p.Wait()
}

// 上传单个图片任务
func UploadURL(info taskData) (bool, string) {
	baseURL := viper.GetString("base_url")
	token := viper.GetString("token")
	if baseURL == "" || token == "" {
		model.DB.Model(&model.ZPurls{}).Where("id = ?", info.ID).Update("status", model.UploadStatusFailed)
		return false, "缺少base_url或token配置"
	}
	apiURL := baseURL + "/api/v3/upload"

	model.DB.Model(&model.ZPurls{}).Where("id = ?", info.ID).Update("status", model.UploadStatusUploading)

	re, filePath := DownloadURL(info.OriginURL)

	// 无论如何，最终删除临时文件
	defer delTempFile(filePath)

	if !re {
		model.DB.Model(&model.ZPurls{}).Where("id = ?", info.ID).Updates(map[string]interface{}{
			"status": model.UploadStatusFailed,
		})
		// 写入日志
		helper.WriteLog("URL下载失败，ID：" + strconv.Itoa(int(info.ID)) + "，错误信息：" + filePath)
		return false, "下载失败"
	}

	params := fmt.Sprintf(`{"dedup":true,"album_id":%d,"watermark":false,"compress":false}`, info.AlbumID)

	client := resty.New().SetTimeout(60 * time.Second)
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("User-Agent", "Zpic Client ").
		SetFile("file", filePath).
		SetFormData(map[string]string{
			"params": params,
		}).
		Post(apiURL)
	if err != nil {
		model.DB.Model(&model.ZPurls{}).Where("id = ?", info.ID).Update("status", model.UploadStatusFailed)
		// 写入日志
		helper.WriteLog("URL上传失败，ID：" + strconv.Itoa(int(info.ID)) + "，错误信息：" + err.Error())
		return false, "上传失败：" + err.Error()
	}

	body := resp.String()
	code := gjson.Get(body, "code").Int()
	if code != 200 {
		msg := gjson.Get(body, "msg").String()
		model.DB.Model(&model.ZPurls{}).Where("id = ?", info.ID).Update("status", model.UploadStatusFailed)
		// 写入日志
		helper.WriteLog("URL上传失败，ID：" + strconv.Itoa(int(info.ID)) + "，错误信息：" + msg)
		return false, "上传失败：" + msg
	}

	model.DB.Model(&model.ZPurls{}).Where("id = ?", info.ID).Updates(map[string]interface{}{
		"imgid":        gjson.Get(body, "data.imgid").String(),
		"url":          gjson.Get(body, "data.url").String(),
		"filename":     gjson.Get(body, "data.filename").String(),
		"image_width":  int(gjson.Get(body, "data.width").Int()),
		"image_height": int(gjson.Get(body, "data.height").Int()),
		"status":       model.UploadStatusSuccess,
	})

	return true, "上传成功"
}

// 删除临时文件
func delTempFile(filePath string) {
	// 检查文件是否存在
	if _, err := os.Stat(filePath); err == nil {
		// 文件存在，删除文件
		err := os.Remove(filePath)
		if err != nil {
			// 写入日志
			helper.WriteLog("删除临时文件失败：" + filePath + "，错误信息：" + err.Error())
		}
	}
}

// 现在图片到本地
func DownloadURL(urlStr string) (bool, string) {
	httpProxy := viper.GetString("http_proxy")
	tempDir := "data/temp/"
	client := resty.New().SetTimeout(60 * time.Second)
	if httpProxy != "" {
		client.SetProxy(httpProxy)
	}

	headResp, err := client.R().
		SetHeader("Referer", urlStr).
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36").
		Head(urlStr)
	if err != nil {
		// 写入日志
		helper.WriteLog(urlStr + "获取文件头失败：" + err.Error())
		return false, "获取文件头失败：" + err.Error()
	}

	contentType := strings.ToLower(headResp.Header().Get("Content-Type"))
	allowedTypes := map[string]string{
		"image/jpeg": "jpg",
		"image/jpg":  "jpg",
		"image/png":  "png",
		"image/bmp":  "bmp",
		"image/gif":  "gif",
		"image/webp": "webp",
	}
	ext, ok := allowedTypes[contentType]
	if !ok {
		// 写入日志
		helper.WriteLog("不支持的文件类型：" + contentType + "，URL：" + urlStr)
		return false, "不支持的文件类型：" + contentType
	}

	contentLength := headResp.Header().Get("Content-Length")
	if contentLength == "" {
		helper.WriteLog(urlStr + "无法获取文件大小")
		return false, "无法获取文件大小"
	}

	size, err := strconv.ParseInt(contentLength, 10, 64)
	if err != nil {
		helper.WriteLog(urlStr + "解析文件大小失败：" + err.Error())
		return false, "解析文件大小失败：" + err.Error()
	}
	if size > 10*1024*1024 {
		return false, "文件超过10MB限制"
	}

	// if err := os.MkdirAll(tempDir, 0o755); err != nil {
	// 	return false, "创建临时目录失败：" + err.Error()
	// }

	fileName := ""
	if parsedURL, err := url.Parse(urlStr); err == nil {
		baseName := path.Base(parsedURL.Path)
		if baseName != "" && baseName != "." && path.Ext(baseName) != "" {
			fileName = baseName
		}
	}
	if fileName == "" {
		fileName = time.Now().Format("20060102150405") + "." + ext
	}

	filePath := path.Join(tempDir, fileName)
	_, err = client.R().
		SetHeader("Referer", urlStr).
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36").
		SetOutput(filePath).
		Get(urlStr)
	if err != nil {
		helper.WriteLog(urlStr + "下载文件失败：" + err.Error())
		return false, "下载失败：" + err.Error()
	}

	return true, filePath
}

// 根据ID批量三次
func (ac *AppCore) DeleteUrlsByIds(data IdsForm) ResData {
	ids := data.Ids

	// 批量删除
	result := model.DB.Where("id IN ?", ids).Delete(&model.ZPurls{})
	if result.Error != nil {
		return ResData{
			Status: false,
			Msg:    "删除URL失败：" + result.Error.Error(),
			Data:   nil,
		}
	}

	return ResData{
		Status: true,
		Msg:    "删除URL成功",
		Data:   nil,
	}
}

// 声明一个函数，函数接受字符串，判断是否是有效的URL格式，然后返回BOOL
func IsValidURL(url string) bool {
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return true
	}
	return false
}
