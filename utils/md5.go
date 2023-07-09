package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
)

func Md5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	tempStr := h.Sum(nil)

	return hex.EncodeToString(tempStr)
}

// 大寫
func MD5Encode(data string) string {
	return strings.ToUpper(Md5Encode(data))
}

// 加密
func MakePassword(plainpwd, salt string) string {
	fmt.Println(111111, salt, plainpwd, salt+plainpwd)
	return Md5Encode(plainpwd + salt)
}

// 加密
func ValidPassword(plainpwd, salt string, password string) bool {
	return Md5Encode(plainpwd+salt) == password
}
