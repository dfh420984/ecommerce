package utils

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UploadFile(c *gin.Context, dir string) (string, error) {
	file, err := c.FormFile("file")
	if err != nil {
		return "", err
	}

	ext := "." + GetFileExt(file.Filename)
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	dst := fmt.Sprintf("%s/%s", dir, filename)

	if err := c.SaveUploadedFile(file, dst); err != nil {
		return "", err
	}

	return "/uploads/" + filename, nil
}

func UploadFiles(c *gin.Context, dir string) ([]string, error) {
	form, err := c.MultipartForm()
	if err != nil {
		return nil, err
	}

	files := form.File["files"]
	if len(files) == 0 {
		files = form.File["file"]
	}

	var paths []string
	for _, file := range files {
		ext := "." + GetFileExt(file.Filename)
		filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
		dst := fmt.Sprintf("%s/%s", dir, filename)

		if err := c.SaveUploadedFile(file, dst); err != nil {
			return nil, err
		}
		paths = append(paths, "/uploads/"+filename)
	}

	return paths, nil
}

func GetFileExt(filename string) string {
	idx := len(filename) - 1
	for idx > 0 {
		if filename[idx] == '.' {
			return filename[idx+1:]
		}
		idx--
	}
	return "jpg"
}

func GetClientIP(c *gin.Context) string {
	ip := c.GetHeader("X-Forwarded-For")
	if ip == "" {
		ip = c.GetHeader("X-Real-IP")
	}
	if ip == "" {
		ip = c.ClientIP()
	}
	return ip
}

func GetUserAgent(c *gin.Context) string {
	return c.GetHeader("User-Agent")
}

func NowDate() string {
	return time.Now().Format("2006-01-02")
}

func NowDateTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
