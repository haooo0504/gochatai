package service

import (
	"encoding/hex"
	"fmt"
	"gochatai/models"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetPostList
// @Security ApiKeyAuth
// @Summary 貼文列表
// @Tags 貼文資料
// @param id query string false "id"
// @Success 200 {string} json{"code","message"}
// @Router /post/getPostList [get]
func GetPostList(c *gin.Context) {
	id := c.Query("id")
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		// handle error
	}
	data := make([]*models.PostWithLikes, 10)
	data = models.GetPostList(uint(userID))

	c.JSON(http.StatusOK, gin.H{
		"code":    0, // 0 成功 -1失敗
		"message": "獲取資料成功",
		"data":    data,
	})
}

// CreatePost
// @Security ApiKeyAuth
// @Summary 創建貼文
// @Tags 貼文資料
// @param author formData string false "用戶名"
// @param title formData string false "標題"
// @param content formData string false "內容"
// @param image formData file false "照片"
// @Success 200 {string} json{"code","message"}
// @Router /post/createPost [post]
func CreatePost(c *gin.Context) {
	post := models.PostInfo{}
	// 取得 Multipart 表單
	form, _ := c.MultipartForm()
	// 從表單中取得貼文的標題和內容
	post.Author = form.Value["author"][0]
	post.Title = form.Value["title"][0]
	post.Content = form.Value["content"][0]

	// 從表單中取得上傳的圖片
	file, _ := c.FormFile("image")
	// 保存圖片到本地文件系統
	// Generate a random string for filename
	buf := make([]byte, 16) // 16 bytes will give us 32 hex characters
	if _, err := rand.Read(buf); err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("generate random string err: %s", err.Error()))
		return
	}
	randomString := hex.EncodeToString(buf)

	filename := fmt.Sprintf("%s-%s", randomString, filepath.Base(file.Filename))
	directory := "assets/images/"
	os.MkdirAll(directory, os.ModePerm) // 確保目錄存在
	path := filepath.Join(directory, filename)
	if err := c.SaveUploadedFile(file, path); err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("upload file err: %s", err.Error()))
		return
	}
	post.ImageURL = "/assets/images/" + filename

	newPost, err := models.CreatePost(&post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, // 0 成功 -1失敗
			"message": "無法創建貼文"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0, // 0 成功 -1失敗
		"message": "創建貼文成功",
		"data":    newPost,
	})
}
