package service

import (
	"fmt"
	"gochatai/models"
	"gochatai/utils"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type UserInput struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// GetUserList
// @Security ApiKeyAuth
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
// @param UserInput body UserInput true "用戶名和密碼"
// @Success 200 {string} json{"code","message"}
// @Router /user/findUserByNameAndPwd [post]
func FindUserByNameAndPwd(c *gin.Context) {
	data := models.UserBasic{}
	var userInput UserInput
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	name := userInput.Name
	password := userInput.Password
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

	// Create the JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = name
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	t, err := token.SignedString([]byte("your_secret_key"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not sign the token"})
		return
	}

	pwd := utils.MakePassword(password, user.Salt)
	data = models.FindUserByNameAndPwd(name, pwd, t)

	c.JSON(http.StatusOK, gin.H{
		"code":    0, // 0 成功 -1失敗
		"message": "登入成功",
		"data":    data,
		"token":   t,
	})
}

// CreateUser
// @Summary 新增用戶
// @Tags 用戶資料
// @param email formData string true "電子郵件"
// @param name formData string true "用戶名"
// @param password formData string true "密碼"
// @param repassword formData string true "確認密碼"
// @Success 200 {string} json{"code","message"}
// @Router /user/createUser [post]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	user.Email = c.PostForm("email")
	user.Name = c.PostForm("name")
	password := c.PostForm("password")
	repassword := c.PostForm("repassword")

	isValidEmail := govalidator.IsEmail(user.Email)
	// Check if the email is empty
	if !isValidEmail {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "Email格式不正確",
			"data":    user,
		})
		return
	}
	if user.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "Email是必填項目",
			"data":    user,
		})
		return
	}
	if len(user.Name) < 1 {
		c.JSON(-1, gin.H{
			"code":    -1, // 0 成功 -1失敗
			"message": "用戶名必須至少有一個字符",
			"data":    user,
		})
		return
	}
	if len(password) < 6 || len(password) > 12 {
		c.JSON(-1, gin.H{
			"code":    -1, // 0 成功 -1失敗
			"message": "密碼長度必須在6到12個字符之間",
			"data":    user,
		})
		return
	}

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
	// Create the JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = user.Name
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	t, err := token.SignedString([]byte("your_secret_key"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not sign the token"})
		return
	}

	// user.Password = password
	user.Password = utils.MakePassword(password, salt)
	user.Salt = salt
	user.Identity = t
	newUser, err := models.CreateUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0, // 0 成功 -1失敗
		"message": "新增用戶成功",
		"data":    newUser,
	})
}

// DeleteUser
// @Security ApiKeyAuth
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
// @Security ApiKeyAuth
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

// GoogleSignIn
// @Summary google登入
// @Tags 用戶資料
// @param idToken formData string true "idToken"
// @Success 200 {string} json{"code","message"}
// @Router /user/googleSignIn [post]
func GoogleSignIn(c *gin.Context) {
	idToken := c.PostForm("idToken")
	user := models.UserBasic{}
	payload, err := utils.ValidateGoogleIdToken(idToken)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    -1,
			"message": "Invalid Google ID token",
		})
		return
	}

	// Create the JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = user.Name
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	t, err := token.SignedString([]byte("your_secret_key"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not sign the token"})
		return
	}

	user.Identity = t

	// Replace these lines with the real values from the ID token
	email := payload.Claims["email"].(string)
	name, _ := utils.GetNameFromIdToken(idToken)
	hasUser := models.FindUserByEmail(email)
	fmt.Println(hasUser)
	if hasUser.Name == "" {
		// User does not exist yet, create a new one
		user.Email = email
		user.Name = name
		newUser, err := models.CreateUser(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
			return
		}
		fmt.Println(user)
		c.JSON(http.StatusOK, gin.H{
			"code":    0, // 0 成功 -1失敗
			"message": "註冊成功",
			"data":    newUser,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    0, // 0 成功 -1失敗
			"message": "登入成功",
			"data":    hasUser,
		})
	}

}
