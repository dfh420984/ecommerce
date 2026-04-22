package handlers

import (
	"shop_api/utils"

	"github.com/gin-gonic/gin"
)

// UploadImage 上传图片
func UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.Fail(c, "请选择要上传的文件")
		return
	}

	// 验证文件类型
	ext := utils.GetFileExt(file.Filename)
	allowedExts := map[string]bool{
		"jpg": true, "jpeg": true, "png": true, "gif": true, "webp": true,
	}
	if !allowedExts[ext] {
		utils.Fail(c, "只支持 jpg、jpeg、png、gif、webp 格式的图片")
		return
	}

	// 验证文件大小（5MB）
	if file.Size > 5*1024*1024 {
		utils.Fail(c, "图片大小不能超过 5MB")
		return
	}

	// 上传文件
	path, err := utils.UploadFile(c, "./uploads")
	if err != nil {
		utils.Fail(c, "上传失败: "+err.Error())
		return
	}

	utils.Success(c, gin.H{
		"url": path,
	})
}

// UploadVideo 上传视频
func UploadVideo(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.Fail(c, "请选择要上传的视频文件")
		return
	}

	// 验证文件类型
	ext := utils.GetFileExt(file.Filename)
	allowedExts := map[string]bool{
		"mp4": true, "webm": true, "ogg": true, "avi": true, "mov": true, "wmv": true, "flv": true,
	}
	if !allowedExts[ext] {
		utils.Fail(c, "只支持 mp4、webm、ogg、avi、mov、wmv、flv 格式的视频")
		return
	}

	// 验证文件大小（50MB）
	if file.Size > 50*1024*1024 {
		utils.Fail(c, "视频大小不能超过 50MB")
		return
	}

	// 上传文件
	path, err := utils.UploadFile(c, "./uploads")
	if err != nil {
		utils.Fail(c, "上传失败: "+err.Error())
		return
	}

	utils.Success(c, gin.H{
		"url": path,
	})
}
