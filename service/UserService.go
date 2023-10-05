package service

import (
	"encoding/hex"
	"fmt"
	"gochatai/models"
	"gochatai/utils"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
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

// // GetUserList
// // @Security ApiKeyAuth
// // @Tags 用戶資料
// // @Success 200 {string} json{"code","message"}
// // @Router /user/getUserList [get]
// func GetUserList(c *gin.Context) {
// 	data := make([]*models.UserBasic, 10)
// 	data = models.GetUserList()

// 	c.JSON(http.StatusOK, gin.H{
// 		"code":    0, // 0 成功 -1失敗
// 		"message": "獲取資料成功",
// 		"data":    data,
// 	})
// }

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
			"message": "此用戶名稱已被他人使用",
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

	user.IsGoogle = false
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
// @param oldname formData string false "舊用戶名"
// @param newname formData string false "新用戶名"
// @param oldpassword formData string false "舊密碼"
// @param newpassword formData string false "新密碼"
// @param phone formData string false "phone"
// @param email formData string false "email"
// @param image formData file false "照片"
// @Success 200 {string} json{"code","message"}
// @Router /user/updateUser [post]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	user.Name = c.PostForm("newname")
	oldpassword := c.PostForm("oldpassword")
	newpassword := c.PostForm("newpassword")
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")

	oldUser := models.FindUserById(user.ID)
	if !oldUser.IsGoogle && !oldUser.IsApple {
		if oldUser.Name == "" {
			c.JSON(400, gin.H{
				"code":    -1, // 0 成功 -1失敗
				"message": "該用戶不存在",
			})
			return
		}

		flag := utils.ValidPassword(oldpassword, oldUser.Salt, oldUser.Password)
		if !flag {
			c.JSON(http.StatusOK, gin.H{
				"code":    -1, // 0 成功 -1失敗
				"message": "密碼錯誤",
			})
			return
		}
	}

	if user.Name != "" {
		data := models.FindUserByName(user.Name)
		if data.Name != "" {
			c.JSON(-1, gin.H{
				"code":    -1, // 0 成功 -1失敗
				"message": "此用戶名稱已被他人使用",
				"data":    user,
			})
			return
		}
	}
	if newpassword != "" {
		salt := fmt.Sprintf("%06d", rand.Int31())
		user.Password = utils.MakePassword(newpassword, salt)
		user.Salt = salt
	}

	// 從表單中取得上傳的圖片
	file, _ := c.FormFile("image")
	if file != nil {
		// 保存圖片到本地文件系統
		// Generate a random string for filename
		buf := make([]byte, 16) // 16 bytes will give us 32 hex characters
		if _, err := rand.Read(buf); err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("generate random string err: %s", err.Error()))
			return
		}
		randomString := hex.EncodeToString(buf)

		filename := fmt.Sprintf("%s-%s", randomString, filepath.Base(file.Filename))
		directory := "assets/userImages/"
		os.MkdirAll(directory, os.ModePerm) // 確保目錄存在
		path := filepath.Join(directory, filename)
		if err := c.SaveUploadedFile(file, path); err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}
		user.ImageURL = "/assets/userImages/" + filename
	}

	_, err := govalidator.ValidateStruct(user)
	curUser, err := models.UpdateUser(user)
	if err != nil {
		c.JSON(400, gin.H{
			"code":    -1,
			"message": "修改用戶失敗",
			"data":    err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "修改用戶成功",
			"data":    curUser,
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

	user.IsGoogle = true

	// Replace these lines with the real values from the ID token
	email := payload.Claims["email"].(string)
	sub := payload.Claims["sub"].(string)
	name, _ := utils.GetNameFromIdToken(idToken)
	hasUser := models.FindUserByGoogleSignIn(email, sub)
	if hasUser.Name == "" {
		// User does not exist yet, create a new one
		user.Email = email
		user.Name = name
		user.TokenSub = sub
		newUser, err := models.CreateUser(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": -1, // 0 成功 -1失敗
				"message": "無法創建用戶"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    0, // 0 成功 -1失敗
			"message": "註冊成功",
			"data":    newUser,
		})
	} else {
		// 用户已经存在，返回存在的用户
		user.ID = hasUser.ID
		curUser, err := models.UpdateUser(user)
		if err != nil {
			c.JSON(400, gin.H{
				"code":    -1,
				"message": "登入失敗",
				"data":    err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code":    0,
				"message": "登入成功",
				"data":    curUser,
			})
		}
	}

}

// RefreshToken
// @Tags 更新token
// @param id formData uint true "id"
// @param name formData string true "用戶名"
// @param token formData string true "token"
// @Success 200 {string} json{"code","message"}
// @Router /user/RefreshToken [post]
func RefreshToken(c *gin.Context) {
	idStr := c.PostForm("id")
	name := c.PostForm("name")
	oldToken := c.PostForm("token")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	// Create the JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	t, err := token.SignedString([]byte("your_secret_key"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not sign the token"})
		return
	}

	data, correct := models.RefreshToken(uint(id), name, oldToken, t)
	if correct {
		c.JSON(http.StatusOK, gin.H{
			"code":    0, // 0 成功 -1失敗
			"message": "更新token成功",
			"data":    data,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1, // 0 成功 -1失敗
			"message": "更新tokeng失敗",
		})
	}

}

// AppleSignIn
// @Summary AppleSignIn
// @Tags 用戶資料
// @param idToken formData string true "idToken"
// @param userName formData string false "userName"
// @Success 200 {string} json{"code","message"}
// @Router /user/appleSignIn [post]
func AppleSignIn(c *gin.Context) {
	idToken := c.PostForm("idToken") // 从请求中获取 Apple ID token
	userName := c.PostForm("userName")
	user := models.UserBasic{}

	// 验证 Apple ID token
	token, err := utils.VerifyAppleToken(idToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    -1,
			"message": "Invalid Apple ID token",
		})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var userEmail string                  // 用于存储邮箱的变量
		email, ok := claims["email"].(string) // 假设email在 claims 中
		sub, ok := claims["sub"].(string)     // 假设email在 claims 中
		if !ok {
			userEmail, ok = claims["sub"].(string)
			if !ok {
				c.JSON(http.StatusBadRequest, gin.H{"code": -1, "message": "Email is required"})
				return
			}
		} else {
			userEmail = email
		}
		userEmail += "AppleId"
		hasUser := models.FindUserByAppleSignIn(userEmail, sub)
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
		user.IsApple = true

		if hasUser.Name == "" {

			// 用户尚未存在，创建一个新的用户
			user.Email = userEmail
			if userName != "" {
				user.Name = userName
			} else {
				user.Name = "user" + fmt.Sprintf("%d", time.Now().Unix())
			}
			user.TokenSub = sub
			// 添加任何其他需要的用户属性

			newUser, err := models.CreateUser(&user)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    -1,
					"message": "無法創建用戶",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"code":    0,
				"message": "註冊成功",
				"data":    newUser,
			})
		} else {
			// 用户已经存在，返回存在的用户
			user.ID = hasUser.ID
			curUser, err := models.UpdateUser(user)
			if err != nil {
				c.JSON(400, gin.H{
					"code":    -1,
					"message": "登入失敗",
					"data":    err.Error(),
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"code":    0,
					"message": "登入成功",
					"data":    curUser,
				})
			}
			// c.JSON(http.StatusOK, gin.H{
			// 	"code":    0,
			// 	"message": "登入成功",
			// 	"data":    hasUser,
			// })
		}
	}
}
