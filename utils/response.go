package utils

import "github.com/gin-gonic/gin"

type Response struct {
	Code    int    `json:"code"`
	Data    any    `json:"data"`
	Message string `json:"message"`
}

func getObj(code int, data any, msg string) gin.H {
	return gin.H{
		"code":    code,
		"data":    data,
		"message": msg,
	}
}
func res(c *gin.Context, code int, data any, msg string) {
	c.JSON(code, getObj(code, data, msg))
}
func Ok(c *gin.Context, data any) {
	res(c, 200, data, "success")
}
func Ok2(c *gin.Context, msg string) {
	res(c, 200, gin.H{}, msg)
}
func Fail(c *gin.Context, code int, msg string) {
	res(c, code, gin.H{}, msg)
}
