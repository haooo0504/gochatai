package models

import (
	"fmt"
	"gochatai/utils"

	"gorm.io/gorm"
)

type Like struct {
	gorm.Model
	UserID uint
	PostID uint
}

func (table *Like) TableName() string {
	return "likes"
}

// AddLike 用于添加点赞
func AddLike(userID uint, postID uint) (int64, error) {
	if UserAlreadyLiked(userID, postID) {
		// 用户已经点过赞了，不再添加
		return 0, fmt.Errorf("user has already liked this post")
	}

	// 创建新的点赞
	like := Like{
		UserID: userID,
		PostID: postID,
	}

	if err := utils.DB.Create(&like).Error; err != nil {
		return 0, err
	}

	// 获取最新的点赞数
	var count int64
	utils.DB.Model(&Like{}).Where("post_id = ?", postID).Count(&count)

	return count, nil
}

// RemoveLike 用于移除点赞
func RemoveLike(userID uint, postID uint) (int64, error) {
	if !UserAlreadyLiked(userID, postID) {
		// 用户没有点过赞，无法移除
		return 0, fmt.Errorf("user has not liked this post")
	}

	utils.DB.Where("user_id = ? AND post_id = ?", userID, postID).Delete(&Like{})

	// 获取最新的点赞数
	var count int64
	utils.DB.Model(&Like{}).Where("post_id = ?", postID).Count(&count)
	fmt.Println(count)
	return count, nil
}

// 用戶按過的所有"讚"
func GetPostWithLikes(postID uint) (PostInfo, error) {
	var post PostInfo
	err := utils.DB.Preload("Likes").Find(&post, postID).Error
	return post, err
}

// UserAlreadyLiked 用於檢查某個用戶是否已經對某個貼文點過讚
func UserAlreadyLiked(userID uint, postID uint) bool {
	var like Like
	if err := utils.DB.Where("user_id = ? AND post_id = ?", userID, postID).First(&like).Error; err != nil {
		// 如果出現錯誤，表示該用戶還未對該貼文點過讚
		return false
	}
	// 如果未出現錯誤，表示該用戶已經對該貼文點過讚
	return true
}

func UserExists(userId uint) bool {
	var count int64
	utils.DB.Model(&UserBasic{}).Where("id = ?", userId).Count(&count)
	return count > 0
}

func PostExists(postId uint) bool {
	var count int64
	utils.DB.Model(&PostInfo{}).Where("id = ?", postId).Count(&count)
	return count > 0
}
