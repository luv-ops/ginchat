package utils

import "golang.org/x/crypto/bcrypt"

// 生成hash密码的计算轮次,2**cost次,判断是否相等也会由这个cost性能影响
const cost = 11

func Encode(originPassword string) (string, error) {
	data, err := bcrypt.GenerateFromPassword([]byte(originPassword), cost)
	return string(data), err
}

func Verify(hash string, input string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(input)) == nil
}
