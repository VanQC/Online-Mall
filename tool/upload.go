package tool

import (
	"github.com/gin-gonic/gin"
	"log"
	"mime/multipart"
	"os"
	"strconv"
)

// todo 这两个方法可以合并
//func UploadToLocalStatic(c *gin.Context, file *multipart.FileHeader, id int, name, basePath string) (filePath string, err error) {
//
//
//}

// UploadProductToLocalStatic 上传到本地文件中
func UploadProductToLocalStatic(c *gin.Context, file *multipart.FileHeader, bossId uint, productName string) (filePath string, err error) {
	bId := strconv.Itoa(int(bossId))
	basePath := "." + "/static/imgs/product/" + "boss" + bId + "/"
	if !DirExistOrNot(basePath) {
		CreateDir(basePath)
	}
	productPath := basePath + productName + ".jpg"
	err = c.SaveUploadedFile(file, productPath)
	if err != nil {
		return "", err
	}
	return "boss" + bId + "/" + productName + ".jpg", err
}

// UploadAvatarToLocalStatic 上传头像
func UploadAvatarToLocalStatic(c *gin.Context, file *multipart.FileHeader, userId uint, userName string) (filePath string, err error) {
	bId := strconv.Itoa(int(userId))
	basePath := "." + "/static/imgs/avatar/" + "user" + bId + "/"
	if !DirExistOrNot(basePath) {
		CreateDir(basePath)
	}
	avatarPath := basePath + userName + ".jpg"
	err = c.SaveUploadedFile(file, avatarPath)
	if err != nil {
		return "", err
	}
	return "user" + bId + "/" + userName + ".jpg", err
}

// DirExistOrNot 判断文件是否存在
func DirExistOrNot(fileAddr string) bool {
	s, err := os.Stat(fileAddr)
	if err != nil {
		log.Println(err)
		return false
	}
	return s.IsDir()
}

// CreateDir 创建文件夹
func CreateDir(dirName string) bool {
	err := os.MkdirAll(dirName, 7550)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
