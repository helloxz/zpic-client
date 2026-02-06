package core

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"zpic-client/helper"
	"zpic-client/model"

	"github.com/zeebo/xxh3"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	MaxFileSize = 10 * 1024 * 1024 // 10MB
)

var (
	allowedExtensions = map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".bmp":  true,
		".gif":  true,
		".webp": true,
	}
)

// ScanListParams 获取扫描任务列表的请求参数
// Page: 当前页码，从1开始
// Limit: 每页显示数量，默认为10
type ScanListParams struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// ScanListResponse 获取扫描任务列表的响应数据结构
// Items: 当前页的任务列表
// Total: 任务总数量
// Page: 当前页码
// Limit: 每页显示数量
type ScanListResponse struct {
	Items []model.ZPtasks `json:"items"`
	Total int64           `json:"total"`
	Page  int             `json:"page"`
	Limit int             `json:"limit"`
}

// GetScanList 获取扫描任务列表
// 支持分页查询，按任务ID降序排列
// params.Page: 页码，默认为1
// params.Limit: 每页数量，默认为10
// 返回任务列表及总数
func (ac *AppCore) GetScanList(params ScanListParams) ResData {
	page := params.Page
	// 页码最小为1
	if page < 1 {
		page = 1
	}
	limit := params.Limit
	// 每页数量最小为1，最大为100
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	// 计算偏移量
	offset := (page - 1) * limit

	var tasks []model.ZPtasks
	// 按ID降序查询，限制返回数量并跳过已显示的数据
	result := model.DB.Order("id desc").Limit(limit).Offset(offset).Find(&tasks)
	if result.Error != nil {
		return ResData{
			Status: false,
			Msg:    "获取任务列表失败：" + result.Error.Error(),
			Data:   nil,
		}
	}

	// 统计任务总数
	var total int64
	model.DB.Model(&model.ZPtasks{}).Count(&total)

	return ResData{
		Status: true,
		Msg:    "获取成功",
		Data: ScanListResponse{
			Items: tasks,
			Total: total,
			Page:  page,
			Limit: limit,
		},
	}
}

// AddScanTaskParams 添加扫描任务的请求参数
// Path: 要扫描的目录路径
type AddScanTaskParams struct {
	Path string `json:"path"`
}

// AddScanTask 添加扫描任务
// 将用户选择的目录路径写入数据库，创建待扫描任务
// 任务初始状态为 PendingScan(0)，成功/失败/总数均为0
// params.Path: 扫描目录的绝对路径
// 返回操作结果
func (ac *AppCore) AddScanTask(params AddScanTaskParams) ResData {
	// 校验目录路径不能为空
	if params.Path == "" {
		return ResData{
			Status: false,
			Msg:    "目录路径不能为空",
			Data:   nil,
		}
	}

	// 创建新任务记录
	task := model.ZPtasks{
		Path:       params.Path,       // 扫描路径
		Status:     model.PendingScan, // 初始状态：待扫描
		SuccessNum: 0,                 // 初始成功数量：0
		FailedNum:  0,                 // 初始失败数量：0
		TotalNum:   0,                 // 初始总数：0
		CreatedAt:  time.Now(),        // 创建时间
		UpdatedAt:  time.Now(),        // 更新时间
	}

	// 写入数据库
	result := model.DB.Create(&task)
	if result.Error != nil {
		return ResData{
			Status: false,
			Msg:    "创建任务失败：" + result.Error.Error(),
			Data:   nil,
		}
	}

	// 异步执行扫描任务，不阻塞AddScanTask
	go ScanTaskURLS(ScanTaskURLSParams{
		TaskID: task.ID,
		Path:   params.Path,
	})

	return ResData{
		Status: true,
		Msg:    "任务创建成功",
		Data:   nil,
	}
}

// ScanTaskURLSParams 扫描任务的参数结构体
// 用于接收任务ID和路径信息
type ScanTaskURLSParams struct {
	TaskID uint   `json:"task_id"`
	Path   string `json:"path"`
}

// ScanTaskURLS 扫描任务目录下的图片文件
// 扫描指定目录中符合条件（指定后缀且大小不超过10MB）的图片文件
// 将扫描结果写入zp_task_urls表，并更新对应任务的扫描状态
// params: 包含任务ID和扫描路径的参数
func ScanTaskURLS(params ScanTaskURLSParams) {
	taskID := params.TaskID
	scanPath := params.Path

	// 检查路径是否存在
	if _, err := os.Stat(scanPath); os.IsNotExist(err) {
		// 路径不存在，将任务状态更新为UploadCompleted
		if updateResult := model.DB.Model(&model.ZPtasks{}).Where("id = ?", taskID).Updates(map[string]interface{}{
			"status":     model.UploadCompleted,
			"updated_at": time.Now(),
		}); updateResult.Error != nil {
			helper.WriteLog("ScanTaskURLS: 更新任务状态失败，任务ID: " + strconv.Itoa(int(taskID)) + "，错误: " + updateResult.Error.Error())
		}
		return
	}

	// 扫描目录下的文件
	fileInfos, err := scanDirectory(scanPath)
	if err != nil {
		helper.WriteLog("ScanTaskURLS: 扫描目录失败，路径: " + scanPath + "，错误: " + err.Error())
		// 更新任务状态为失败
		return
	}

	// 如果没有扫描到文件，直接返回
	if len(fileInfos) == 0 {
		return
	}

	// 使用事务提交数据，保障一致性
	err = model.DB.Transaction(func(tx *gorm.DB) error {
		// 批量插入扫描到的文件信息到zp_task_urls表
		// 使用INSERT OR IGNORE处理重复hash的情况
		var taskUrls []model.ZPTaskUrls
		for _, info := range fileInfos {
			taskUrls = append(taskUrls, model.ZPTaskUrls{
				TaskID:     taskID,
				OriginPath: info.FullPath,
				FileName:   info.FileName,
				FileSize:   info.FileSize,
				FileHash:   info.FileHash,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
				Status:     model.URLPending,
			})
		}

		if len(taskUrls) > 0 {
			// 使用Clauses配置INSERT OR IGNORE，处理重复hash
			if createResult := tx.Clauses(&clause.Insert{
				Modifier: "OR IGNORE",
			}).Create(&taskUrls); createResult.Error != nil {
				return createResult.Error
			}
		}

		// 获取实际插入成功的记录数量
		var insertedCount int64
		if countResult := tx.Model(&model.ZPTaskUrls{}).Where("task_id = ?", taskID).Count(&insertedCount); countResult.Error != nil {
			return countResult.Error
		}

		// 更新任务状态为ScanCompleted=1，并更新TotalNum
		if updateResult := tx.Model(&model.ZPtasks{}).Where("id = ?", taskID).Updates(map[string]interface{}{
			"status":     model.ScanCompleted,
			"total_num":  int(insertedCount),
			"updated_at": time.Now(),
		}); updateResult.Error != nil {
			return updateResult.Error
		}

		return nil
	})

	if err != nil {
		helper.WriteLog("ScanTaskURLS: 事务提交失败，任务ID: " + strconv.Itoa(int(taskID)) + "，错误: " + err.Error())
	}
}

// FileInfo 扫描到的文件信息结构体
// 用于存储扫描到的图片文件的详细信息
type FileInfo struct {
	FullPath string // 文件绝对完整路径
	FileName string // 文件名
	FileSize int64  // 文件大小（字节）
	FileHash string // 文件hash
}

// scanDirectory 扫描指定目录下的图片文件
// scanPath: 要扫描的目录路径
// 返回扫描到的文件信息列表和错误信息
// 扫描规则：
//   - 只扫描指定后缀的文件：.jpg, .jpeg, .png, .bmp, .gif, .webp（不区分大小写）
//   - 文件大小不能超过10MB
//   - 不递归扫描子目录
func scanDirectory(scanPath string) ([]FileInfo, error) {
	var files []FileInfo

	// 读取目录下的所有条目
	entries, err := os.ReadDir(scanPath)
	if err != nil {
		return nil, err
	}

	// 遍历目录条目
	for _, entry := range entries {
		// 只处理文件，不处理子目录
		if entry.IsDir() {
			continue
		}

		// 获取文件名
		fileName := entry.Name()

		// 获取文件扩展名（转小写）
		ext := strings.ToLower(filepath.Ext(fileName))

		// 检查文件扩展名是否在允许列表中
		if !allowedExtensions[ext] {
			continue
		}

		// 获取文件的完整路径
		fullPath := filepath.Join(scanPath, fileName)

		// 获取文件信息
		info, err := entry.Info()
		if err != nil {
			continue
		}

		// 获取文件大小
		fileSize := info.Size()

		// 检查文件大小是否超过10MB
		if fileSize > MaxFileSize {
			continue
		}

		// 获取文件hash
		fileHash, err := getFileHash(fullPath)
		if err != nil {
			continue
		}

		// 添加到文件列表
		files = append(files, FileInfo{
			FullPath: fullPath,
			FileName: fileName,
			FileSize: fileSize,
			FileHash: fileHash,
		})
	}

	return files, nil
}

// getFileHash 计算文件的hash值
// 使用xxh3算法计算文件的哈希值
// filePath: 文件的完整路径
// 返回文件的hash值
func getFileHash(filePath string) (string, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 计算文件hash（xxh3算法）
	hasher := xxh3.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}
	fileHash := fmt.Sprintf("%x", hasher.Sum64())

	return fileHash, nil
}

// hexEncodeToString 将hash值转换为十六进制字符串
// 用于将xxh3计算得到的hash值转换为可存储的字符串格式
func hexEncodeToString(hash [16]byte) string {
	return fmt.Sprintf("%x", hash)
}

// DeleteTasksParams 删除任务的请求参数
// Ids: 要删除的任务ID数组
type DeleteTasksParams struct {
	Ids []uint `json:"ids"`
}

// DeleteTaskResult 删除任务的返回结果
// DeletedCount: 成功删除的任务数量
// SkippedCount: 跳过的任务数量（因为状态为上传中）
type DeleteTaskResult struct {
	DeletedCount int    `json:"deleted_count"`
	SkippedCount int    `json:"skipped_count"`
	Message      string `json:"message"`
}

// DeleteTasks 删除任务
// 批量删除指定的任务及其关联的URL记录
// 只删除status非Uploading(2)的任务，Uploading状态的任务会被跳过
// 使用事务保障数据完整性，同时删除zp_tasks和zp_task_urls表中的记录
// params.Ids: 要删除的任务ID列表
// 返回删除结果详情
func (ac *AppCore) DeleteTasks(params DeleteTasksParams) ResData {
	// 校验至少选择一个任务
	if len(params.Ids) == 0 {
		return ResData{
			Status: false,
			Msg:    "请选择要删除的任务",
			Data:   nil,
		}
	}

	// 筛选出非Uploading状态的任务ID
	// Uploading状态为2，不能删除
	var deletableIds []uint
	var skippedCount int64

	// 查询状态不为Uploading(2)的任务
	queryResult := model.DB.Model(&model.ZPtasks{}).
		Where("id IN ? AND status != ?", params.Ids, model.Uploading).
		Pluck("id", &deletableIds)

	if queryResult.Error != nil {
		return ResData{
			Status: false,
			Msg:    "查询可删除任务失败：" + queryResult.Error.Error(),
			Data:   nil,
		}
	}

	// 计算被跳过的任务数量（总选择数 - 可删除数）
	skippedCount = int64(len(params.Ids)) - int64(len(deletableIds))

	// 如果没有可删除的任务，直接返回
	if len(deletableIds) == 0 {
		msg := "所选任务均为上传中状态，无法删除"
		if skippedCount > 0 {
			msg = fmt.Sprintf("%d个任务均为上传中状态，无法删除", skippedCount)
		}
		return ResData{
			Status: false,
			Msg:    msg,
			Data: DeleteTaskResult{
				DeletedCount: 0,
				SkippedCount: int(skippedCount),
				Message:      msg,
			},
		}
	}

	// 使用事务确保数据完整性
	// 1. 先删除zp_task_urls表中关联的URL记录
	// 2. 再删除zp_tasks表中的任务记录
	err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 批量删除关联的URL记录
		// 使用WHERE ... IN ...一次性删除所有匹配TaskID的记录，避免循环删除
		if deleteUrlsResult := tx.Where("task_id IN ?", deletableIds).Delete(&model.ZPTaskUrls{}); deleteUrlsResult.Error != nil {
			return fmt.Errorf("删除关联URL失败：%w", deleteUrlsResult.Error)
		}

		// 批量删除任务记录
		if deleteTasksResult := tx.Where("id IN ?", deletableIds).Delete(&model.ZPtasks{}); deleteTasksResult.Error != nil {
			return fmt.Errorf("删除任务失败：%w", deleteTasksResult.Error)
		}

		return nil
	})

	if err != nil {
		return ResData{
			Status: false,
			Msg:    "删除失败：" + err.Error(),
			Data: DeleteTaskResult{
				DeletedCount: 0,
				SkippedCount: int(skippedCount),
				Message:      err.Error(),
			},
		}
	}

	// 构建返回消息
	message := fmt.Sprintf("成功删除%d个任务", len(deletableIds))
	if skippedCount > 0 {
		message = fmt.Sprintf("成功删除%d个任务，跳过%d个上传中任务", len(deletableIds), skippedCount)
	}

	return ResData{
		Status: true,
		Msg:    message,
		Data: DeleteTaskResult{
			DeletedCount: len(deletableIds),
			SkippedCount: int(skippedCount),
			Message:      message,
		},
	}
}

// SelectScanDirectory 选择扫描目录
// 打开系统目录选择对话框，让用户选择要扫描的文件夹
// 返回用户选择的目录路径，用户取消返回空字符串
func (ac *AppCore) SelectScanDirectory() (string, error) {
	// 打开目录选择对话框
	path, err := runtime.OpenDirectoryDialog(appCtx, runtime.OpenDialogOptions{
		Title: "选择要扫描的目录",
	})
	if err != nil {
		return "", err
	}
	return path, nil
}

// GetTotalPages 获取扫描任务的总页数
// 根据任务总数和每页10条计算总页数
// 返回总页数，最少返回1
func (ac *AppCore) GetTotalPages() int {
	var total int64
	model.DB.Model(&model.ZPtasks{}).Count(&total)
	limit := 10 // 与GetScanList保持一致
	// 无数据时返回1
	if total <= 0 {
		return 1
	}
	// 计算总页数
	pages := int(total) / limit
	// 如果有余数，总页数加1
	if int(total)%limit > 0 {
		pages++
	}
	return pages
}

// GetScanTaskCount 获取扫描任务的总数
// 统计 zp_tasks 表中的所有记录数量
// 返回任务总数
func (ac *AppCore) GetScanTaskCount() int64 {
	var count int64
	model.DB.Model(&model.ZPtasks{}).Count(&count)
	return count
}

// TaskURLCounts 任务URL统计结果结构体
// 用于存储单个任务的URL成功和失败数量统计
type TaskURLCounts struct {
	SuccessCount int64 // 成功上传的URL数量
	FailedCount  int64 // 上传失败的URL数量
	TotalCount   int64 // 成功+失败的总数
}

// UpdateOneTask 更新单个上传任务的进度
// 1. 查询status=1或2的任务（扫描完成/上传中），按id增序排列，取一条
// 2. 获取该任务的ID和TotalNum
// 3. 查询该任务关联的所有URL的统计信息（成功数、失败数）
// 4. 判断是否完成：成功+失败 == TotalNum
//   - 完成：更新status为UploadCompleted(3)，同时更新SuccessNum和FailedNum
//   - 未完成但有进度（成功数或失败数>0）：更新status为Uploading(2)，更新SuccessNum和FailedNum
//   - 未完成且无进度：只更新SuccessNum和FailedNum，状态不变
//
// 性能优化：使用单次查询统计URL数量，减少数据库访问
func UpdateOneTask() (bool, error) {
	// Step 1: 查询status=1(扫描完成)或status=2(上传中)的任务，按id增序排序，只取一条
	var task model.ZPtasks
	queryResult := model.DB.Where("status IN ?", []int8{model.ScanCompleted, model.Uploading}).Order("id asc").First(&task)
	if queryResult.Error != nil {
		if queryResult.Error == gorm.ErrRecordNotFound {
			// 没有需要处理的任务，返回false
			return false, nil
		}
		return false, queryResult.Error
	}

	// 没有查询到任务
	if task.ID == 0 {
		return false, nil
	}

	// Step 2 & 3: 使用单次查询统计该任务下所有URL的成功和失败数量
	// 使用原生SQL进行聚合查询，性能优于多次查询
	var counts TaskURLCounts
	countQuery := `
		SELECT
			COALESCE(SUM(CASE WHEN status = ? THEN 1 ELSE 0 END), 0) as success_count,
			COALESCE(SUM(CASE WHEN status = ? THEN 1 ELSE 0 END), 0) as failed_count,
			COALESCE(SUM(CASE WHEN status IN (?, ?) THEN 1 ELSE 0 END), 0) as total_count
		FROM %s
		WHERE task_id = ?
	`
	countSQL := fmt.Sprintf(countQuery, model.ZPTaskUrls{}.TableName(), model.URLSuccess, model.URLFailed, model.URLSuccess, model.URLFailed)

	countResult := model.DB.Raw(countSQL, task.ID).Scan(&counts)
	if countResult.Error != nil {
		return false, countResult.Error
	}

	// Step 4: 判断任务状态并更新
	isCompleted := counts.TotalCount == int64(task.TotalNum)
	hasProgress := counts.SuccessCount > 0 || counts.FailedCount > 0

	// 构建更新数据
	updateData := map[string]interface{}{
		"success_num": counts.SuccessCount,
		"failed_num":  counts.FailedCount,
		"updated_at":  time.Now(),
	}

	// 根据完成状态和进度更新任务状态
	if isCompleted {
		// 任务完成：更新为UploadCompleted(3)
		updateData["status"] = model.UploadCompleted
	} else if hasProgress {
		// 任务进行中（有进度但未完成）：更新为Uploading(2)
		updateData["status"] = model.Uploading
	}
	// 未完成且无进度：只更新数量，状态不变

	// 执行更新
	updateResult := model.DB.Model(&model.ZPtasks{}).Where("id = ?", task.ID).Updates(updateData)
	if updateResult.Error != nil {
		return false, updateResult.Error
	}

	return true, nil
}

// BatchTaskURLCounts 批量任务URL统计结果
// 用于批量查询多个任务的URL统计信息
type BatchTaskURLCounts struct {
	TaskID       uint  `gorm:"column:task_id"`
	SuccessCount int64 `gorm:"column:success_count"`
	FailedCount  int64 `gorm:"column:failed_count"`
	TotalCount   int64 `gorm:"column:total_count"`
}

// UpdateOneTaskBatch 批量更新上传任务进度
// 性能优化：使用批量子查询，将N+1问题优化为固定3次数据库操作
// batchSize: 每次处理的任务数量，默认为10
func UpdateOneTaskBatch(batchSize int) (int, error) {
	if batchSize <= 0 {
		batchSize = 10
	}

	// Step 1: 查询所有status=1(扫描完成)或status=2(上传中)的任务，按id增序排序
	var tasks []model.ZPtasks
	queryResult := model.DB.Where("status IN ?", []int8{model.ScanCompleted, model.Uploading}).Order("id asc").Limit(batchSize).Find(&tasks)
	if queryResult.Error != nil {
		return 0, queryResult.Error
	}

	// 没有需要处理的任务
	if len(tasks) == 0 {
		return 0, nil
	}

	// 收集所有任务ID
	taskIDs := make([]uint, len(tasks))
	for i, task := range tasks {
		taskIDs[i] = task.ID
	}

	// Step 2: 批量统计所有任务的URL数量（单次查询，使用GROUP BY）
	var batchCounts []BatchTaskURLCounts
	batchQuery := fmt.Sprintf(`
		SELECT
			task_id,
			COALESCE(SUM(CASE WHEN status = %d THEN 1 ELSE 0 END), 0) as success_count,
			COALESCE(SUM(CASE WHEN status = %d THEN 1 ELSE 0 END), 0) as failed_count,
			COALESCE(SUM(CASE WHEN status IN (%d, %d) THEN 1 ELSE 0 END), 0) as total_count
		FROM %s
		WHERE task_id IN ?
		GROUP BY task_id
	`, model.URLSuccess, model.URLFailed, model.URLSuccess, model.URLFailed, model.ZPTaskUrls{}.TableName())

	countResult := model.DB.Raw(batchQuery, taskIDs).Scan(&batchCounts)
	if countResult.Error != nil {
		return 0, countResult.Error
	}

	// 构建taskID到统计结果的映射，方便后续使用
	countMap := make(map[uint]BatchTaskURLCounts, len(batchCounts))
	for _, bc := range batchCounts {
		countMap[bc.TaskID] = bc
	}

	// Step 3: 批量更新所有任务状态（使用CASE WHEN）
	updatedAt := time.Now()

	// 构建批量更新数据
	type updateData struct {
		TaskID      uint
		SuccessNum  int64
		FailedNum   int64
		NewStatus   int8
		NeedsStatus bool // 是否需要更新状态
	}

	var updates []updateData
	for _, task := range tasks {
		counts, exists := countMap[task.ID]
		if !exists {
			// 如果没有查询到统计数据，跳过
			continue
		}

		isCompleted := counts.TotalCount == int64(task.TotalNum)
		hasProgress := counts.SuccessCount > 0 || counts.FailedCount > 0

		newStatus := task.Status
		needsStatus := false

		if isCompleted {
			newStatus = model.UploadCompleted
			needsStatus = true
		} else if hasProgress && task.Status != model.Uploading {
			newStatus = model.Uploading
			needsStatus = true
		}

		if int64(task.SuccessNum) != counts.SuccessCount || int64(task.FailedNum) != counts.FailedCount || needsStatus {
			updates = append(updates, updateData{
				TaskID:      task.ID,
				SuccessNum:  counts.SuccessCount,
				FailedNum:   counts.FailedCount,
				NewStatus:   newStatus,
				NeedsStatus: needsStatus,
			})
		}
	}

	// 执行批量更新
	if len(updates) > 0 {
		// 使用事务确保数据一致性
		err := model.DB.Transaction(func(tx *gorm.DB) error {
			for _, u := range updates {
				updateData := map[string]interface{}{
					"success_num": u.SuccessNum,
					"failed_num":  u.FailedNum,
					"updated_at":  updatedAt,
				}
				if u.NeedsStatus {
					updateData["status"] = u.NewStatus
				}
				if result := tx.Model(&model.ZPtasks{}).Where("id = ?", u.TaskID).Updates(updateData); result.Error != nil {
					return result.Error
				}
			}
			return nil
		})
		if err != nil {
			return 0, err
		}
	}

	return len(updates), nil
}
