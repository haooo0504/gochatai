package models

import (
	"fmt"
	"gochatai/utils"

	"gorm.io/gorm"
)

type PostInfo struct {
	gorm.Model
	Author   string `gorm:"type:varchar(100)"`
	Title    string `gorm:"type:varchar(100)"`
	Content  string `gorm:"type:text"`
	ImageURL string `gorm:"type:varchar(255)"`
}

func (table *PostInfo) TableName() string {
	return "post_info"
}

// 獲取貼文列表
func GetPostList() []*PostInfo {
	data := make([]*PostInfo, 10)
	utils.DB.Find(&data)
	for _, v := range data {
		fmt.Println(v)
	}
	return data
}

// 創建貼文
func CreatePost(post PostInfo) *gorm.DB {
	// 建立新的貼文
	return utils.DB.Create(&post)
}
