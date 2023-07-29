package models

import (
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

// AddLike adds a like for a specific post by a user
func AddLike(userID uint, postID uint) *gorm.DB {
	like := Like{
		UserID: userID,
		PostID: postID,
	}
	return utils.DB.Create(&like)
}

// RemoveLike removes a like for a specific post by a user
func RemoveLike(userID, postID uint) *gorm.DB {
	like := Like{
		UserID: userID,
		PostID: postID,
	}
	return utils.DB.Delete(&like)
}

// 用戶按過的所有"讚"
func GetPostWithLikes(postID uint) (PostInfo, error) {
	var post PostInfo
	err := utils.DB.Preload("Likes").Find(&post, postID).Error
	return post, err
}
