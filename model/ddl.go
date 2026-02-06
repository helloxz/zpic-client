package model

import "time"

// UploadStatus constants
const (
	UploadStatusPending   = 0 // 未开始
	UploadStatusUploading = 1 // 上传中
	UploadStatusSuccess   = 2 // 已完成
	UploadStatusFailed    = 3 // 上传失败
)

// URL上传表结构
type ZPurls struct {
	ID          uint      `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	OriginURL   string    `gorm:"column:origin_url;not null" json:"origin_url"`
	ImgID       string    `gorm:"column:imgid;size:16" json:"imgid"`
	AlbumID     int       `gorm:"column:album_id;default:0;not null" json:"album_id"`
	FileName    string    `gorm:"column:filename" json:"filename"`
	URL         string    `gorm:"type:text;column:url" json:"url"`
	ImageWidth  int       `gorm:"column:image_width" json:"image_width"`
	ImageHeight int       `gorm:"column:image_height" json:"image_height"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
	Status      int8      `gorm:"column:status;default:0;not null;index:idx_status" json:"status"`
}

// 自定义表名
func (ZPurls) TableName() string {
	return "zp_upload_urls"
}

// 任务表status状态
const (
	PendingScan = 0 // 待扫描
	// 扫描完成
	ScanCompleted = 1
	// 上传中
	Uploading = 2
	// 上传完成
	UploadCompleted = 3
)

// 任务表接口
type ZPtasks struct {
	ID         uint      `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	Path       string    `gorm:"column:path" json:"path"`
	SuccessNum int       `gorm:"column:success_num;default:0;not null" json:"success_num"`             // 成功数量
	FailedNum  int       `gorm:"column:failed_num;default:0;not null" json:"failed_num"`               // 失败数量
	TotalNum   int       `gorm:"column:total_num;default:0;not null" json:"total_num"`                 // 总数量
	Status     int8      `gorm:"column:status;default:0;not null;index:idx_task_status" json:"status"` // 任务状态
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (ZPtasks) TableName() string {
	return "zp_tasks"
}

// 任务表里面的URL状态
const (
	URLPending   = 0 // 待上传
	URLUploading = 1 // 上传中
	URLSuccess   = 2 // 成功
	URLFailed    = 3 // 失败
)

// 任务表里面的URL
type ZPTaskUrls struct {
	ID     uint `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	TaskID uint `gorm:"column:task_id;not null;index:idx_task_id" json:"task_id"` // 任务ID
	// 原始文件路径
	OriginPath string `gorm:"column:origin_path;not null" json:"origin_path"`
	// 临时文件路径，允许为空
	TempPath string `gorm:"column:temp_path" json:"temp_path"`
	// 后端返回的图片ID
	Imgid string `gorm:"column:imgid" json:"imgid"`
	// 上传后得到的URL，允许为空
	URL         string `gorm:"type:text;column:url" json:"url"`
	ImageWidth  int    `gorm:"column:image_width"`
	ImageHeight int    `gorm:"column:image_height"`
	// 文件名，允许空
	FileName string `gorm:"column:filename" json:"filename"`
	// 原始文件大小，字节单位，不能为空
	FileSize int64 `gorm:"column:file_size;not null" json:"file_size"`
	// 真实文件大小，字节单位，默认为0
	RealFileSize int64 `gorm:"column:real_file_size;default:0;not null" json:"real_file_size"`
	// 文件Hash，不能为空，且唯一
	FileHash  string    `gorm:"column:file_hash;not null;uniqueIndex" json:"file_hash"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	Status    int8      `gorm:"column:status;default:0;not null;index:idx_task_url_status" json:"status"` // 任务状态
}

func (ZPTaskUrls) TableName() string {
	return "zp_task_urls"
}
