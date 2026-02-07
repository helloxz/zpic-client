package core

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
	"zpic-client/model"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// ExportTaskURLS 导出任务URL列表到CSV文件
// @param taskID 任务ID，前端传递过来
// @return ResData 返回操作结果
func (ac *AppCore) ExportTaskURLS(taskID uint) ResData {
	// 参数验证
	if taskID == 0 {
		return ResData{Status: false, Msg: "任务ID不能为空"}
	}

	// 检查任务是否存在
	var task model.ZPtasks
	result := model.DB.First(&task, taskID)
	if result.Error != nil {
		return ResData{Status: false, Msg: "任务不存在"}
	}

	// 任务未完成时不允许导出
	if task.Status != model.UploadCompleted {
		return ResData{Status: false, Msg: "任务未上传完成，无法导出"}
	}

	// 查询该任务下所有成功的URL记录
	var taskUrls []model.ZPTaskUrls
	queryResult := model.DB.Where("task_id = ? AND status = ?", taskID, model.URLSuccess).Find(&taskUrls)
	if queryResult.Error != nil {
		return ResData{Status: false, Msg: "查询数据失败：" + queryResult.Error.Error()}
	}

	// 没有可导出的数据
	if len(taskUrls) == 0 {
		return ResData{Status: false, Msg: "没有可导出的数据"}
	}

	// 打开保存文件对话框
	defaultFilename := fmt.Sprintf("task_%d_export_%s.csv", taskID, time.Now().Format("2006-01-02"))
	filePath, err := runtime.SaveFileDialog(appCtx, runtime.SaveDialogOptions{
		Title:           "导出URL列表",
		DefaultFilename: defaultFilename,
		Filters: []runtime.FileFilter{
			{DisplayName: "CSV Files", Pattern: "*.csv"},
		},
	})

	if err != nil {
		return ResData{Status: false, Msg: "打开保存对话框失败：" + err.Error()}
	}

	// 用户取消选择
	if filePath == "" {
		return ResData{Status: false, Msg: "已取消导出"}
	}

	// 创建CSV文件
	file, err := os.Create(filePath)
	if err != nil {
		return ResData{Status: false, Msg: "创建文件失败：" + err.Error()}
	}
	defer file.Close()

	// 写入UTF-8 BOM，防止Excel打开中文乱码
	file.WriteString("\xEF\xBB\xBF")

	// 创建CSV写入器
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入CSV表头
	header := []string{
		"图片ID",
		"文件名",
		"链接",
		"宽",
		"高",
		"文件大小",
		"上传时间",
	}
	if err := writer.Write(header); err != nil {
		return ResData{Status: false, Msg: "写入表头失败：" + err.Error()}
	}

	// 写入数据行
	for _, item := range taskUrls {
		record := []string{
			item.Imgid,
			item.FileName,
			item.URL,
			strconv.Itoa(item.ImageWidth),
			strconv.Itoa(item.ImageHeight),
			strconv.FormatInt(item.RealFileSize, 10),
			item.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		if err := writer.Write(record); err != nil {
			return ResData{Status: false, Msg: "写入数据失败：" + err.Error()}
		}
	}

	return ResData{
		Status: true,
		Msg:    fmt.Sprintf("成功导出 %d 条数据", len(taskUrls)),
		Data:   nil,
	}
}
