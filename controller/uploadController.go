package controller

import (
	"GinChat/service"
	"GinChat/utils"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// UploadFile
// @Summary 图片文件上传
// @Description 上传聊天图片，返回可访问URL
// @Tags 文件相关
// @Accept multipart/form-data
// @Param img formData file true "图片文件(jpg/png/gif)"
// @Success 200 {object} utils.Response{} "上传成功"
// @Failure 400 {object} utils.Response{} "参数错误/文件非法"
// @Failure 500 {object} utils.Response{} "服务器错误"
// @Router /upload/img [post]
func UploadFile(c *gin.Context) {

	file, err := c.FormFile("file")
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}
	// 限制大小 5MB
	const maxSize = 5 << 20
	if file.Size > maxSize {
		utils.Fail(c, 400, "文件不能超过5MB")
		return
	}
	// 校验后缀
	ext := filepath.Ext(file.Filename)
	allowExt := map[string]bool{".jpg": true, ".png": true, ".jpeg": true, ".gif": true}
	if !allowExt[ext] {
		utils.Fail(c, 400, "文件格式错误")
		return
	}
	uploadFilePath, err := service.UploadFile(file, c)
	if err != nil {
		utils.Fail(c, 500, "服务器错误")
		return
	}
	utils.Ok(c, viper.GetString("file.domain")+uploadFilePath)
}
