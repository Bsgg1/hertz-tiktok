package db

import (
	"tiktok/pkg/constants"
	"tiktok/pkg/errno"
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

func AddNewComment(comment *Comment) error {
	ok, _ := CheckUserExistById(comment.UserId)
	if !ok {
		return errno.UserIsNotExistErr
	}
	ok, _ = CheckVideoExistById(comment.VideoId)
	if !ok {
		return errno.VideoIsNotExistErr
	}
	err := DB.Create(comment).Error
	if err != nil {
		return err
	}
	return nil
}

func DeleteCommentById(comment_id int64) error {
	if ok, _ := CheckCommentExist(comment_id); !ok {
		return errno.CommentIsNotExistErr
	}
	comment := &Comment{}
	err := DB.Where("id = ?", comment_id).Delete(comment).Error
	if err != nil {
		return err
	}
	return nil
}

func CheckCommentExist(comment_id int64) (bool, error) {
	comment := &Comment{}
	err := DB.Where("id = ?", comment_id).Find(comment).Error
	if err != nil {
		return false, err
	}
	if comment.ID == 0 {
		return false, nil
	}
	return true, nil
}

func GetCommentListByVideoID(video_id int64) ([]*Comment, error) {
	var comment_list []*Comment
	if ok, _ := CheckVideoExistById(video_id); !ok {
		return comment_list, errno.VideoIsNotExistErr
	}
	err := DB.Table(constants.CommentTableName).Where("video_id = ?", video_id).Find(&comment_list).Error
	if err != nil {
		return comment_list, err
	}
	return comment_list, nil
}

func GetCommentCountByVideoID(video_id int64) (int64, error) {
	var sum int64
	err := DB.Model(&Comment{}).Where("video_id = ?", video_id).Count(&sum).Error
	if err != nil {
		return sum, err
	}
	return sum, nil
}
