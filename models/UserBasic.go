package models

import (
	"fmt"
	"gochatai/utils"

	"gorm.io/gorm"
)

type UserBasic struct {
	gorm.Model
	Name          string `gorm:"not null"`
	Password      string
	Phone         string
	Email         string `gorm:"not null;valid:\"email\""`
	Identity      string
	ClientIp      string
	ClientPort    string
	Salt          string
	LoginTime     uint64
	HeartbeatTime uint64
	LogOutTime    uint64
	IsLogout      bool
	DeviceInfo    string
	CanUseTime    uint64
	ImageURL      string
	Likes         []Like `gorm:"foreignKey:UserID"`
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

func GetUserList() []*UserBasic {
	data := make([]*UserBasic, 10)
	utils.DB.Find(&data)
	for _, v := range data {
		fmt.Println(v)
	}
	return data
}

func FindUserByNameAndPwd(name string, password string, token string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name = ? and password = ?", name, password).First(&user)
	// token加密
	// str := fmt.Sprintf("%d", time.Now().Unix())
	// temp := utils.MD5Encode(str)
	utils.DB.Model(&user).Where("id = ?", user.ID).Update("identity", token)
	return user
}

func FindUserByName(name string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name = ?", name).First(&user)
	return user
}

func FindUserByPhone(phone string) *gorm.DB {
	user := UserBasic{}
	return utils.DB.Where("phone = ?", phone).First(&user)
}

func FindUserByEmail(email string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("email = ?", email).First(&user)
	return user
}

func CreateUser(user *UserBasic) (*UserBasic, error) {
	result := utils.DB.Create(user)
	if result.Error != nil {
		// 如果创建操作失败，返回nil和错误
		return nil, result.Error
	}
	// 如果创建操作成功，返回新创建的用户和nil错误
	return user, nil
}

func DeleteUser(user UserBasic) *gorm.DB {
	return utils.DB.Delete(&user)
}

func UpdateUser(user UserBasic) *gorm.DB {
	return utils.DB.Model(&user).Updates(UserBasic{Name: user.Name, Password: user.Password, Phone: user.Phone, Email: user.Email, ImageURL: user.ImageURL, Salt: user.Salt})
}
