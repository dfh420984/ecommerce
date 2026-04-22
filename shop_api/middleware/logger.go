package middleware

import (
	"bytes"
	"context"
	"io"
	"shop_api/database"
	"shop_api/models"
	"shop_api/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// isWriteMethod 判断是否为写操作
func isWriteMethod(method string) bool {
	method = strings.ToUpper(method)
	return method == "POST" || method == "PUT" || method == "DELETE" || method == "PATCH"
}

// OperationLog 操作日志中间件 - 只记录写操作
func OperationLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 只记录写操作
		if !isWriteMethod(c.Request.Method) {
			c.Next()
			return
		}

		// 读取请求体
		var requestBody string
		if c.Request.Body != nil {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			requestBody = string(bodyBytes)
			// 重置请求体，以便后续 handler 可以读取
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// 记录开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 获取用户ID
		userID := utils.GetUserID(c)
		if userID == 0 {
			return
		}

		// 获取用户信息
		var user models.User
		if err := database.GetDB().First(&user, userID).Error; err != nil {
			return
		}

		// 计算处理时间
		latency := time.Since(startTime)

		// 创建操作日志
		log := models.OperationLog{
			UserID:    userID,
			Username:  user.Username,
			Action:    c.Request.Method + " " + c.Request.URL.Path,
			Target:    c.Param("id"),
			Content:   requestBody,
			IP:        utils.GetClientIP(c),
			UserAgent: utils.GetUserAgent(c),
		}

		// 异步写入数据库，避免影响响应速度
		go func() {
			database.GetDB().Create(&log)
		}()

		// 在控制台输出日志
		utils.Info("[OperationLog] User: %s, Action: %s, Target: %s, IP: %s, Latency: %v, Status: %d",
			user.Username, log.Action, log.Target, log.IP, latency, c.Writer.Status())
	}
}

func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := utils.GetClientIP(c)
		key := "rate_limit:" + ip

		count, err := database.GetRedis().Incr(context.Background(), key).Result()
		if err != nil {
			c.Next()
			return
		}

		if count == 1 {
			database.GetRedis().Expire(context.Background(), key, time.Minute)
		}

		if count > 100 {
			utils.FailWithCode(c, 429, "too many requests")
			c.Abort()
			return
		}

		c.Next()
	}
}
