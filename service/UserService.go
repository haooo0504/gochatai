package service

import (
	"fmt"
	"gochatai/models"
	"gochatai/utils"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

// GetUserList
// @Tags 用戶資料
// @Success 200 {string} json{"code","message"}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	data := make([]*models.UserBasic, 10)
	data = models.GetUserList()

	c.JSON(http.StatusOK, gin.H{
		"code":    0, // 0 成功 -1失敗
		"message": "獲取資料成功",
		"data":    data,
	})
}

// GetUserList
// @Summary 用戶登入
// @Tags 用戶資料
// @param name query string false "用戶名"
// @param password query string false "密碼"
// @Success 200 {string} json{"code","message"}
// @Router /user/findUserByNameAndPwd [get]
func FindUserByNameAndPwd(c *gin.Context) {
	data := models.UserBasic{}

	name := c.Query("name")
	password := c.Query("password")
	user := models.FindUserByName(name)
	if user.Name == "" {
		c.JSON(400, gin.H{
			"code":    -1, // 0 成功 -1失敗
			"message": "該用戶不存在",
		})
		return
	}

	flag := utils.ValidPassword(password, user.Salt, user.Password)
	if !flag {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1, // 0 成功 -1失敗
			"message": "密碼錯誤",
		})
		return
	}
	pwd := utils.MakePassword(password, user.Salt)
	data = models.FindUserByNameAndPwd(name, pwd)
	c.JSON(http.StatusOK, gin.H{
		"code":    0, // 0 成功 -1失敗
		"message": "登入成功",
		"data":    data,
	})
}

// CreateUser
// @Summary 新增用戶
// @Tags 用戶資料
// @param name query string false "用戶名"
// @param password query string false "密碼"
// @param repassword query string false "確認密碼"
// @Success 200 {string} json{"code","message"}
// @Router /user/createUser [get]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	user.Name = c.Query("name")
	password := c.Query("password")
	repassword := c.Query("repassword")

	salt := fmt.Sprintf("%06d", rand.Int31())

	data := models.FindUserByName(user.Name)
	if data.Name != "" {
		c.JSON(-1, gin.H{
			"code":    -1, // 0 成功 -1失敗
			"message": "此用戶名稱已註冊",
			"data":    user,
		})
		return
	}
	if password != repassword {
		c.JSON(-1, gin.H{
			"code":    -1, // 0 成功 -1失敗
			"message": "兩次密碼不同",
			"data":    user,
		})
		return
	}
	// user.Password = password
	user.Password = utils.MakePassword(password, salt)
	user.Salt = salt
	models.CreateUser(user)
	c.JSON(http.StatusOK, gin.H{
		"code":    0, // 0 成功 -1失敗
		"message": "新增用戶成功",
		"data":    user,
	})
}

// DeleteUser
// @Summary 刪除用戶
// @Tags 用戶資料
// @param id query string false "id"
// @Success 200 {string} json{"code","message"}
// @Router /user/deleteUser [get]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.Query("id"))
	user.ID = uint(id)
	models.DeleteUser(user)
	c.JSON(http.StatusOK, gin.H{
		"code":    0, // 0 成功 -1失敗
		"message": "刪除用戶成功",
		"data":    user,
	})
}

// UpdateUser
// @Summary 修改用戶
// @Tags 用戶資料
// @param id formData string false "id"
// @param name formData string false "用戶名"
// @param password formData string false "密碼"
// @param phone formData string false "phone"
// @param email formData string false "email"
// @Success 200 {string} json{"code","message"}
// @Router /user/updateUser [post]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	user.Name = c.PostForm("name")
	user.Password = c.PostForm("password")
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")

	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{
			"code":    -1, // 0 成功 -1失敗
			"message": "修改用戶失敗",
			"data":    user,
		})
	} else {
		models.UpdateUser(user)
		c.JSON(http.StatusOK, gin.H{
			"code":    0, // 0 成功 -1失敗
			"message": "修改用戶成功",
			"data":    user,
		})
	}

}