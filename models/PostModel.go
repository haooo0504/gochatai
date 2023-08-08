package models

import (
	"gochatai/utils"
	"time"

	"gorm.io/gorm"
)

type PostInfo struct {
	gorm.Model
	Author    string `gorm:"type:varchar(100)"`
	AuthorId  string
	AuthorImg string
	Title     string `gorm:"type:varchar(100)"`
	Content   string `gorm:"type:text"`
	ImageURL  string `gorm:"type:varchar(255)"`
	Likes     []Like `gorm:"foreignKey:PostID"`
}

func (table *PostInfo) TableName() string {
	return "post_info"
}

type PostWithLikes struct {
	PostInfo
	LikeCount int
	UserLiked bool
}

// 獲取五天內的貼文列表
func GetPostList(userID uint) []*PostWithLikes {
	data := make([]*PostInfo, 0)

	// 获取五天前的日期
	fiveDaysAgo := time.Now().AddDate(0, 0, -5)

	// 只获取日期在五天内的帖子
	utils.DB.Table("post_info").
		Joins("left join user_basic on post_info.author = user_basic.name").
		Preload("Likes").
		Select("post_info.*, user_basic.image_url as author_img").
		Where("post_info.created_at > ?", fiveDaysAgo).
		Order("post_info.created_at desc"). // 从新到旧排序
		Find(&data)

	result := make([]*PostWithLikes, 0)
	for _, v := range data {
		userLiked := false
		for _, like := range v.Likes {
			if like.UserID == userID {
				userLiked = true
				break
			}
		}
		postWithLikes := &PostWithLikes{
			PostInfo:  *v,
			LikeCount: len(v.Likes),
			UserLiked: userLiked,
		}
		result = append(result, postWithLikes)
	}
	return result
}

// 創建貼文
func CreatePost(post *PostInfo) (*PostInfo, error) {
	// 建立新的貼文
	result := utils.DB.Create(post)
	if result.Error != nil {
		// 如果创建操作失败，返回nil和错误
		return nil, result.Error
	}
	// 如果创建操作成功，返回新创建的貼文和nil错误
	return post, nil
}

// 刪除貼文
func DeletePost(post *PostInfo) *gorm.DB {
	println(post)
	result := utils.DB.Delete(post)

	if result.Error != nil {
		println(result.Error)
		return result
	}
	return result
}
