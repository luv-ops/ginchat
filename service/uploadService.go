package service

import (
	"crypto/rand"
	"encoding/hex"
	"mime/multipart"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

const (
	UPLOAD_PATH = "./static/upload"
	RESPATH     = "/static/upload"
)

func UploadFile(file *multipart.FileHeader, c *gin.Context) (string, error) {
	ext := filepath.Ext(file.Filename)
	randBytes := make([]byte, 16)
	_, _ = rand.Read(randBytes)
	fileName := hex.EncodeToString(randBytes) + ext
	savePath := UPLOAD_PATH + "/" + fileName
	// 4. 保存文件到本地
	err := c.SaveUploadedFile(file, savePath)
	if err != nil {
		return "", err
	}
	savePath = RESPATH + "/" + fileName
	return savePath, nil
}
