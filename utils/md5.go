package utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
)

// 生成随机盐值
func generateSalt() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
func Md5Encode(password string) (string, string) {
	hash := md5.New()
	salt := generateSalt()
	hash.Write([]byte(password + salt))
	tempStr := hash.Sum(nil)
	return hex.EncodeToString(tempStr), salt
}
func selfUseEncode(password string) string {
	hash := md5.New()
	hash.Write([]byte(password))
	tempStr := hash.Sum(nil)
	return hex.EncodeToString(tempStr)
}

func Equal(password string, salt string, storePassword string) bool {
	encode := selfUseEncode(password + salt)
	return encode == storePassword
}
