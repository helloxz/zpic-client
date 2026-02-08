package core

import (
	"fmt"
	"io"
	"mime"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	"zpic-client/helper"
	"zpic-client/model"
)

// 这个文件主要是扫描上传任务的批量上传功能

// uploadResult 上传结果结构体
// 用于在并发处理中传递每个URL的上传结果
type uploadResult struct {
	ID           uint
	Success      bool
	ImgID        string
	URL          string
	ImageWidth   int
	ImageHeight  int
	RealFileSize int64
	TempFilePath string
	Error        string
}

// BatchUpload 批量上传任务
// 1. 查询40条待上传的URL记录（status=0），按id增序
// 2. 批量更新状态为上传中（status=1）
// 3. 使用5个并发处理文件
// 4. 等待所有处理完成后，使用事务批量更新数据库
// 5. 清理临时文件
func BatchUpload() {
	var urls []model.ZPTaskUrls
	result := model.DB.Where("status = ?", model.URLPending).Order("id asc").Limit(40).Find(&urls)
	if result.Error != nil {
		helper.WriteLog("BatchUpload: 查询待上传记录失败：" + result.Error.Error())
		return
	}

	if len(urls) == 0 {
		return
	}

	var ids []uint
	for _, url := range urls {
		ids = append(ids, url.ID)
	}

	model.DB.Model(&model.ZPTaskUrls{}).Where("id IN ?", ids).Update("status", model.URLUploading)

	results := make(chan uploadResult, len(urls))
	for _, url := range urls {
		u := url
		go func() {
			defer func() {
				if r := recover(); r != nil {
					results <- uploadResult{ID: u.ID, Success: false, Error: fmt.Sprintf("panic: %v", r)}
				}
			}()
			r := processAndUpload(u)
			results <- r
		}()
	}

	var tempFiles []string
	for i := 0; i < len(urls); i++ {
		r := <-results
		uploadResults = append(uploadResults, r)
		if r.TempFilePath != "" {
			tempFiles = append(tempFiles, r.TempFilePath)
		}
	}
	close(results)

	batchUpdateResults()
	cleanupTempFiles(tempFiles)
}

var uploadResults []uploadResult

// batchUpdateResults 批量更新上传结果
func batchUpdateResults() {
	updatedCount := 0
	for _, r := range uploadResults {
		updateData := map[string]interface{}{
			"updated_at": time.Now(),
		}

		if r.Success {
			updateData["imgid"] = r.ImgID
			updateData["url"] = r.URL
			updateData["image_width"] = r.ImageWidth
			updateData["image_height"] = r.ImageHeight
			updateData["real_file_size"] = r.RealFileSize
			updateData["status"] = model.URLSuccess
		} else {
			updateData["status"] = model.URLFailed
		}

		result := model.DB.Model(&model.ZPTaskUrls{}).Where("id = ?", r.ID).Updates(updateData)
		if result.Error != nil {
			helper.WriteLog(fmt.Sprintf("BatchUpload: 更新失败 ID=%d, err=%v", r.ID, result.Error))
		} else if result.RowsAffected > 0 {
			updatedCount++
		}
	}

	uploadResults = nil
}

// processAndUpload 处理单个URL记录并上传
// 返回上传结果
func processAndUpload(url model.ZPTaskUrls) uploadResult {
	filePath := url.OriginPath

	// if _, err := os.Stat(filePath); os.IsNotExist(err) {
	// 	filePath = url.TempPath
	// }

	exists := true
	if filePath == "" {
		exists = false
	} else if _, err := os.Stat(filePath); os.IsNotExist(err) {
		exists = false
	}

	if !exists {
		return uploadResult{
			ID:    url.ID,
			Error: "文件不存在",
		}
	}

	processedPath, err := processFile(filePath, url.FileName)
	if err != nil {
		return uploadResult{
			ID:    url.ID,
			Error: err.Error(),
		}
	}

	req := UploadReq{
		FilePath: processedPath,
		AlbumID:  0,
	}
	success, uploadResp := UploadZpic(req)

	if success {
		return uploadResult{
			ID:           url.ID,
			Success:      true,
			ImgID:        uploadResp.ImgID,
			URL:          uploadResp.URL,
			ImageWidth:   int(uploadResp.Width),
			ImageHeight:  int(uploadResp.Height),
			RealFileSize: uploadResp.Size,
			TempFilePath: processedPath,
		}
	}

	return uploadResult{
		ID:           url.ID,
		Success:      false,
		TempFilePath: processedPath,
	}
}

// processFile 根据OS和MIME类型处理文件
// 返回处理后的文件路径
func processFile(filePath string, fileName string) (string, error) {
	ext := strings.ToLower(filepath.Ext(fileName))
	mimeType := mime.TypeByExtension(ext)

	// 获取临时目录的绝对路径，避免程序目录权限问题
	destDir := getTempDir()
	// if err := os.MkdirAll(destDir, 0755); err != nil {
	// 	return "", fmt.Errorf("创建临时目录失败: %w", err)
	// }
	destPath := filepath.Join(destDir, fileName)

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return "", fmt.Errorf("获取文件信息失败: %w", err)
	}
	// 文件小于300KB直接复制，无需压缩
	if fileInfo.Size() < 300*1024 {
		return copyFile(filePath, destPath)
	}

	switch mimeType {
	case "image/jpeg":
		return compressJpeg(filePath, destPath)
	case "image/png":
		return compressPng(filePath, destPath)
	default:
		return copyFile(filePath, destPath)
	}
}

// compressJpeg 压缩JPEG图片
func compressJpeg(srcPath string, destPath string) (string, error) {
	binDir := getBinDir()

	optimizePath := filepath.Join(binDir, "jpegoptim")
	if runtime.GOOS == "windows" {
		optimizePath += ".exe"
	}

	destDir := filepath.Dir(destPath)

	cmd := exec.Command(optimizePath, "--max=80", "-s", "--all-progressive", srcPath, "-d", destDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		helper.WriteLog(fmt.Sprintf("jpegoptim压缩失败: err=%v, output=%s", err, string(output)))
		return copyFile(srcPath, destPath)
	}

	return destPath, nil
}

// compressPng 压缩PNG图片
func compressPng(srcPath string, destPath string) (string, error) {
	binDir := getBinDir()

	optimizePath := filepath.Join(binDir, "oxipng")
	if runtime.GOOS == "windows" {
		optimizePath += ".exe"
	}

	cmd := exec.Command(optimizePath,
		"-o", "3",
		"--strip", "safe",
		"-t", "1",
		"--out", destPath,
		srcPath,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		helper.WriteLog(fmt.Sprintf("oxipng压缩失败: err=%v, output=%s", err, string(output)))
		return copyFile(srcPath, destPath)
	}

	return destPath, nil
}

// getBinDir 获取压缩工具所在目录的绝对路径
// 优先使用当前工作目录（wails dev 模式），失败后使用可执行文件目录
func getBinDir() string {
	// 获取运行目录
	cwd := helper.GetRunDir()
	binDir := filepath.Join(cwd, "bin")
	if runtime.GOOS == "windows" {
		return filepath.Join(binDir, "windows")
	} else if runtime.GOOS == "darwin" {
		return filepath.Join(binDir, "darwin")
	} else if runtime.GOOS == "linux" {
		return filepath.Join(binDir, "linux")
	}
	return binDir
}

// getTempDir 获取临时目录的绝对路径

// copyFile 复制文件到目标路径
func copyFile(srcPath string, destPath string) (string, error) {
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		return "", fmt.Errorf("源文件不存在: %s", srcPath)
	}

	srcFile, err := os.Open(srcPath)
	if err != nil {
		return "", fmt.Errorf("打开源文件失败: %w", err)
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("创建目标文件失败: %w", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, srcFile); err != nil {
		return "", fmt.Errorf("复制文件失败: %w", err)
	}

	return destPath, nil
}

// cleanupTempFiles 清理本次批量上传创建的临时文件
// 通过记录处理过程中创建的所有临时文件路径进行精确清理
func cleanupTempFiles(tempFiles []string) {
	tempDir := getTempDir()
	for _, filePath := range tempFiles {
		if filePath == "" || !strings.HasPrefix(filePath, tempDir) {
			continue
		}
		if err := os.Remove(filePath); err != nil {
			helper.WriteLog(fmt.Sprintf("清理临时文件失败: %s, err=%v", filePath, err))
		}
	}
}

// RetryTask 重试失败的上传任务
// 参数 taskID: 任务ID
// 返回 ResData: 包含操作结果（成功/失败）及提示信息
func (ac *AppCore) RetryTask(taskID uint) ResData {
	// 1. 通过任务ID查询 zp_tasks 表
	var task model.ZPtasks
	result := model.DB.First(&task, taskID)
	if result.Error != nil {
		return ResData{
			Status: false,
			Msg:    "任务不存在",
			Data:   nil,
		}
	}

	// 2. 检查任务状态，必须是 UploadCompleted 状态才能重试
	if task.Status != model.UploadCompleted {
		return ResData{
			Status: false,
			Msg:    "进行中的任务不支持重试",
			Data:   nil,
		}
	}

	// 3. 检查是否有失败文件
	if task.FailedNum <= 0 {
		return ResData{
			Status: false,
			Msg:    "没有失败文件，不需要重试",
			Data:   nil,
		}
	}

	// 4. 更新 zp_task_urls 表中对应任务ID且状态为 URLFailed 的记录为 URLPending
	updateResult := model.DB.Model(&model.ZPTaskUrls{}).
		Where("task_id = ? AND status = ?", taskID, model.URLFailed).
		Updates(map[string]interface{}{
			"status":     model.URLPending,
			"updated_at": time.Now(),
		})

	if updateResult.Error != nil {
		return ResData{
			Status: false,
			Msg:    "更新失败文件状态失败：" + updateResult.Error.Error(),
			Data:   nil,
		}
	}

	return ResData{
		Status: true,
		Msg:    fmt.Sprintf("已成功重试 %d 个失败文件", updateResult.RowsAffected),
		Data:   nil,
	}
}
