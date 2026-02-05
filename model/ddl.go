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
