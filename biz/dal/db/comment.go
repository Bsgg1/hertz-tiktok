package db

import (
	"tiktok/pkg/constants"
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID          int64          `json:"id"`
	UserId      int64          `json:"user_id"`
	VideoId     int64          `json:"video_id"`
	CommentText string         `json:"comment_text"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (Comment) TableName() string {
	return constants.CommentTableName
}
