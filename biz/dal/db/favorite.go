package db

import (
	"tiktok/biz/mw/redis"
	"tiktok/pkg/constants"
	"time"

	"gorm.io/gorm"
)

var rdFavorite redis.Favorite

type Favorites struct {
	ID        int64          `json:"id"`
	UserId    int64          `json:"user_id"`
	VideoId   int64          `json:"video_id"`
	CreatedAt time.Time      `json:"create_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"delete_at"`
}

func (Favorites) TableName() string {
	return constants.FavoritesTableName
}

func AddNewFavorite(favorite *Favorites) (bool, error) {
	err := DB.Create(favorite).Error
	if err != nil {
		return false, err
	}
	// 喜欢之后，直接放到缓存
	if rdFavorite.CheckLiked(favorite.VideoId) {
		rdFavorite.AddLiked(favorite.UserId, favorite.VideoId)
	}
	if rdFavorite.CheckLike(favorite.UserId) {
		rdFavorite.AddLike(favorite.UserId, favorite.VideoId)
	}

	return true, nil
}

func DeleteFavorite(favorite *Favorites) (bool, error) {
	err := DB.Where("video_id = ? AND user_id = ?", favorite.VideoId, favorite.UserId).Delete(favorite).Error
	if err != nil {
		return false, err
	}
	if rdFavorite.CheckLiked(favorite.VideoId) {
		rdFavorite.DelLiked(favorite.UserId, favorite.VideoId)
	}
	if rdFavorite.CheckLike(favorite.UserId) {
		rdFavorite.DelLike(favorite.UserId, favorite.VideoId)
	}
	return true, nil
}

func QueryFavoriteExist(user_id, video_id int64) (bool, error) {
	//检查是否存在喜欢的时候，先从缓存里边去查
	if rdFavorite.CheckLiked(video_id) {
		return rdFavorite.ExistLiked(user_id, video_id), nil
	}
	if rdFavorite.CheckLike(user_id) {
		return rdFavorite.ExistLike(user_id, video_id), nil
	}
	var sum int64
	err := DB.Model(&Favorites{}).Where("video_id = ? AND user_id = ?", video_id, user_id).Count(&sum).Error
	if err != nil {
		return false, err
	}
	if sum == 0 {
		return false, nil
	}
	return true, nil
}

func QueryTotalFavoritedByAuthorID(author_id int64) (int64, error) {
	var sum int64
	err := DB.Table(constants.FavoritesTableName).Joins("JOIN videos ON likes.video_id = videos.id").
		Where("videos.author_id = ?", author_id).Count(&sum).Error
	if err != nil {
		return 0, err
	}
	return sum, nil
}

func getFavoriteIdList(user_id int64) ([]int64, error) {
	var favorite_actions []Favorites
	err := DB.Where("user_id = ?", user_id).Find(&favorite_actions).Error
	if err != nil {
		return nil, err
	}
	var result []int64
	for _, v := range favorite_actions {
		result = append(result, v.VideoId)
	}
	return result, nil
}

func GetFavoriteIdList(user_id int64) ([]int64, error) {
	if rdFavorite.CheckLike(user_id) {
		return rdFavorite.GetLike(user_id), nil
	}
	return getFavoriteIdList(user_id)
}

func GetFavoriteCountByUserID(user_id int64) (int64, error) {
	if rdFavorite.CheckLike(user_id) {
		return rdFavorite.CountLike(user_id)
	}
	// 不在缓存里边，查完数据库更新缓存
	videos, err := getFavoriteIdList(user_id)
	if err != nil {
		return 0, err
	}

	go func(user int64, videos []int64) {
		for _, video := range videos {
			rdFavorite.AddLiked(user, video)
		}
	}(user_id, videos)

	return int64(len(videos)), nil
}

func getFavoriterIdList(video_id int64) ([]int64, error) {
	var favorite_actions []Favorites
	err := DB.Where("video_id = ?", video_id).Find(&favorite_actions).Error
	if err != nil {
		return nil, err
	}
	var result []int64
	for _, v := range favorite_actions {
		result = append(result, v.UserId)
	}
	return result, nil
}

func GetFavoriterIdList(video_id int64) ([]int64, error) {
	if rdFavorite.CheckLiked(video_id) {
		return rdFavorite.GetLiked(video_id), nil
	}
	return getFavoriterIdList(video_id)
}

func GetFavoriteCount(video_id int64) (int64, error) {
	if rdFavorite.CheckLiked(video_id) {
		return rdFavorite.CountLiked(video_id)
	}

	likes, err := getFavoriterIdList(video_id)
	if err != nil {
		return 0, err
	}

	go func(users []int64, video int64) {
		for _, user := range users {
			rdFavorite.AddLiked(user, video)
		}
	}(likes, video_id)
	return int64(len(likes)), nil
}
