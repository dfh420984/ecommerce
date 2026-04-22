package handlers

import (
	"shop_api/database"
	"shop_api/models"
	"shop_api/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RegisterInput struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Phone    string `json:"phone"`
}

type AdminRegisterInput struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Nickname string `json:"nickname"`
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Fail(c, "参数错误: "+err.Error())
		return
	}

	var count int64
	database.GetDB().Model(&models.User{}).Where("username = ?", input.Username).Count(&count)
	if count > 0 {
		utils.Fail(c, "用户名已存在")
		return
	}

	hash, err := utils.HashPassword(input.Password)
	if err != nil {
		utils.Fail(c, "密码加密失败")
		return
	}

	user := models.User{
		Username: input.Username,
		Password: hash,
		Phone:    input.Phone,
		Nickname: input.Username,
		Status:   models.UserStatusEnabled,
		UserType: models.UserTypeNormal,
	}

	if err := database.GetDB().Create(&user).Error; err != nil {
		utils.Fail(c, "注册失败")
		return
	}

	utils.Success(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
	})
}

func AdminRegister(c *gin.Context) {
	var input AdminRegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Fail(c, "参数错误: "+err.Error())
		return
	}

	var count int64
	database.GetDB().Model(&models.User{}).Where("username = ?", input.Username).Count(&count)
	if count > 0 {
		utils.Fail(c, "用户名已存在")
		return
	}

	hash, err := utils.HashPassword(input.Password)
	if err != nil {
		utils.Fail(c, "密码加密失败")
		return
	}

	nickname := input.Nickname
	if nickname == "" {
		nickname = input.Username
	}

	user := models.User{
		Username: input.Username,
		Password: hash,
		Nickname: nickname,
		Status:   models.UserStatusEnabled,
		UserType: models.UserTypeAdmin,
	}

	if err := database.GetDB().Create(&user).Error; err != nil {
		utils.Fail(c, "注册失败")
		return
	}

	utils.Success(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"nickname": user.Nickname,
	})
}

func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Fail(c, "参数错误: "+err.Error())
		return
	}

	var user models.User
	if err := database.GetDB().Where("username = ?", input.Username).First(&user).Error; err != nil {
		utils.Fail(c, "用户不存在")
		return
	}

	if !utils.CheckPassword(input.Password, user.Password) {
		utils.Fail(c, "密码错误")
		return
	}

	if user.Status != models.UserStatusEnabled {
		utils.Fail(c, "账号已被禁用")
		return
	}

	token, err := utils.GenerateToken(user.ID, 72)
	if err != nil {
		utils.Fail(c, "生成令牌失败")
		return
	}

	utils.Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"nickname": user.Nickname,
			"avatar":   user.Avatar,
			"phone":    user.Phone,
		},
	})
}

func AdminLogin(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Fail(c, "参数错误: "+err.Error())
		return
	}

	var user models.User
	if err := database.GetDB().Where("username = ?", input.Username).First(&user).Error; err != nil {
		utils.Fail(c, "用户不存在")
		return
	}

	if !utils.CheckPassword(input.Password, user.Password) {
		utils.Fail(c, "密码错误")
		return
	}

	if user.Status != models.UserStatusEnabled {
		utils.Fail(c, "账号已被禁用")
		return
	}

	if user.UserType != models.UserTypeAdmin {
		utils.Fail(c, "非管理员账号")
		return
	}

	token, err := utils.GenerateToken(user.ID, 72)
	if err != nil {
		utils.Fail(c, "生成令牌失败")
		return
	}

	utils.Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"nickname": user.Nickname,
			"avatar":   user.Avatar,
			"phone":    user.Phone,
			"user_type": user.UserType,
		},
	})
}

func GetUserInfo(c *gin.Context) {
	userID := utils.GetUserID(c)

	var user models.User
	if err := database.GetDB().First(&user, userID).Error; err != nil {
		utils.Fail(c, "用户不存在")
		return
	}

	utils.Success(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"nickname": user.Nickname,
		"avatar":   user.Avatar,
		"phone":    user.Phone,
		"email":    user.Email,
	})
}

func UpdateUserInfo(c *gin.Context) {
	userID := utils.GetUserID(c)

	var input struct {
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
		Phone    string `json:"phone"`
		Email    string `json:"email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Fail(c, "参数错误")
		return
	}

	updates := map[string]interface{}{}
	if input.Nickname != "" {
		updates["nickname"] = input.Nickname
	}
	if input.Avatar != "" {
		updates["avatar"] = input.Avatar
	}
	if input.Phone != "" {
		updates["phone"] = input.Phone
	}
	if input.Email != "" {
		updates["email"] = input.Email
	}

	if err := database.GetDB().Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		utils.Fail(c, "更新失败")
		return
	}

	utils.Success(c, nil)
}

func ChangePassword(c *gin.Context) {
	userID := utils.GetUserID(c)

	var input struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Fail(c, "参数错误")
		return
	}

	var user models.User
	if err := database.GetDB().First(&user, userID).Error; err != nil {
		utils.Fail(c, "用户不存在")
		return
	}

	if !utils.CheckPassword(input.OldPassword, user.Password) {
		utils.Fail(c, "原密码错误")
		return
	}

	hash, err := utils.HashPassword(input.NewPassword)
	if err != nil {
		utils.Fail(c, "密码加密失败")
		return
	}

	database.GetDB().Model(&user).Update("password", hash)

	utils.Success(c, nil)
}

func WechatLogin(c *gin.Context) {
	var input struct {
		Code string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Fail(c, "参数错误")
		return
	}

	utils.Success(c, gin.H{
		"token": "mock-token-for-wechat",
		"user": gin.H{
			"id":       1,
			"username": "wechat_user",
			"nickname": "微信用户",
		},
	})
}

// GetUsers 管理员获取用户列表
func GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	userType := c.Query("user_type")
	status := c.Query("status")
	keyword := c.Query("keyword")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	query := database.GetDB().Model(&models.User{})

	// 按用户类型筛选
	if userType != "" {
		query = query.Where("user_type = ?", userType)
	}

	// 按状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 关键词搜索（用户名、昵称、手机号）
	if keyword != "" {
		query = query.Where("username LIKE ? OR nickname LIKE ? OR phone LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	var users []models.User
	if err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&users).Error; err != nil {
		utils.Fail(c, "获取用户列表失败")
		return
	}

	// 隐藏密码
	for i := range users {
		users[i].Password = ""
	}

	utils.PageSuccess(c, users, total, page, pageSize)
}

// GetUser 管理员获取用户详情
func GetUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := database.GetDB().First(&user, id).Error; err != nil {
		utils.Fail(c, "用户不存在")
		return
	}

	// 隐藏密码
	user.Password = ""

	utils.Success(c, user)
}

// UpdateUserStatus 管理员更新用户状态
func UpdateUserStatus(c *gin.Context) {
	id := c.Param("id")

	var input struct {
		Status int8 `json:"status" binding:"required,oneof=0 1"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Fail(c, "参数错误: "+err.Error())
		return
	}

	var user models.User
	if err := database.GetDB().First(&user, id).Error; err != nil {
		utils.Fail(c, "用户不存在")
		return
	}

	database.GetDB().Model(&user).Update("status", input.Status)

	utils.Success(c, nil)
}

// DeleteUser 管理员删除用户
func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := database.GetDB().First(&user, id).Error; err != nil {
		utils.Fail(c, "用户不存在")
		return
	}

	// 不允许删除自己
	currentUserID := utils.GetUserID(c)
	if user.ID == currentUserID {
		utils.Fail(c, "不能删除当前登录账号")
		return
	}

	database.GetDB().Delete(&user)

	utils.Success(c, nil)
}
