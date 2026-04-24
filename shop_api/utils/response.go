package utils

import (
	"fmt"
	"log"
	"net/http"
	"shop_api/models"
	"time"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  "success",
		Data: data,
	})
}

func Fail(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: 400,
		Msg:  msg,
	})
}

func FailWithCode(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
	})
}

func Unauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, Response{
		Code: 401,
		Msg:  "unauthorized",
	})
}

type PageResult struct {
	List       interface{} `json:"list"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

func PageSuccess(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	Success(c, PageResult{
		List:       list,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

// SuccessWithPage 分页成功响应（PageSuccess的别名）
func SuccessWithPage(c *gin.Context, list interface{}, page, pageSize int, total int) {
	PageSuccess(c, list, int64(total), page, pageSize)
}

func GetUserID(c *gin.Context) uint64 {
	v, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	return v.(uint64)
}

func GetUser(c *gin.Context) *models.User {
	v, exists := c.Get("user")
	if !exists {
		return nil
	}
	return v.(*models.User)
}

// Info 输出信息日志
func Info(format string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf(format, args...)
	log.Printf("[INFO] %s - %s", timestamp, msg)
}

// Error 输出错误日志
func Error(format string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf(format, args...)
	log.Printf("[ERROR] %s - %s", timestamp, msg)
}
