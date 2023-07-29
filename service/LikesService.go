package service

import (
	"gochatai/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func HandleLike(c *gin.Context) {
	// 這裡只是一個假設的範例，實際上你需要根據你的應用來取得userID和postID
	userID := c.Query("userID") // 假設userID和postID都在請求的參數中
	postID := c.Query("postID")
	userIDUint, _ := strconv.ParseUint(userID, 10, 64) // 將 string 轉換為 uint
	postIDUint, _ := strconv.ParseUint(postID, 10, 64) // 將 string 轉換為 uint

	// 呼叫AddLike函數
	result := models.AddLike(uint(userIDUint), uint(postIDUint))

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to process like"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Successfully added like"})
	}
}

// AddLike
// @Security ApiKeyAuth
// @Summary 貼文按讚
// @Tags 貼文按讚
// @param userId formData string false "用戶ID"
// @param postId formData string false "貼文ID"
// @Success 200 {string} json{"code","message"}
// @Router /like/addLike [post]
func AddLike(c *gin.Context) {
	like := models.Like{}
	userId, err := strconv.ParseUint(c.PostForm("userId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	postId, err := strconv.ParseUint(c.PostForm("postId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	like.UserID = uint(userId)
	like.PostID = uint(postId)
	models.AddLike(like.UserID, like.PostID)
	c.JSON(http.StatusOK, gin.H{
		"code":    0, // 0 成功 -1失敗
		"message": "按讚成功",
		"data":    like,
	})
}
