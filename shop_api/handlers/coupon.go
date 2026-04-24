package handlers

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"shop_api/database"
	"shop_api/models"
	"shop_api/types"
	"shop_api/utils"
)

// ============ 后台管理接口 ============

// AdminCouponList 获取优惠券列表
func AdminCouponList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")
	name := c.Query("name")

	offset := (page - 1) * pageSize

	query := database.GetDB().Model(&models.Coupon{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	var total int64
	query.Count(&total)

	var coupons []models.Coupon
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&coupons).Error; err != nil {
		utils.Fail(c, "获取优惠券列表失败")
		return
	}

	utils.SuccessWithPage(c, coupons, page, pageSize, int(total))
}

// AdminCouponDetail 获取优惠券详情
func AdminCouponDetail(c *gin.Context) {
	id := c.Param("id")

	var coupon models.Coupon
	if err := database.GetDB().First(&coupon, id).Error; err != nil {
		utils.Fail(c, "优惠券不存在")
		return
	}

	utils.Success(c, coupon)
}

// AdminCreateCoupon 创建优惠券
func AdminCreateCoupon(c *gin.Context) {
	var input struct {
		Name          string  `json:"name" binding:"required"`
		Type          int8    `json:"type" binding:"required,min=1,max=3"`
		DiscountValue float64 `json:"discount_value" binding:"required,min=0"`
		MinAmount     float64 `json:"min_amount"`
		MaxDiscount   float64 `json:"max_discount"`
		TotalCount    int     `json:"total_count"`
		PerUserLimit  int     `json:"per_user_limit" binding:"min=1"`
		ValidType     int8    `json:"valid_type" binding:"required,min=1,max=2"`
		StartTime     *string `json:"start_time"`
		EndTime       *string `json:"end_time"`
		ValidDays     int     `json:"valid_days"`
		Status        int8    `json:"status"`
		IsNewUser     int8    `json:"is_new_user"`
		Description   string  `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Fail(c, "参数错误: "+err.Error())
		return
	}

	// 验证有效期
	if input.ValidType == 1 {
		if input.StartTime == nil || input.EndTime == nil {
			utils.Fail(c, "固定时间类型必须设置开始和结束时间")
			return
		}
	} else if input.ValidType == 2 {
		if input.ValidDays <= 0 {
			utils.Fail(c, "领取后有效天数必须大于0")
			return
		}
	}

	coupon := models.Coupon{
		Name:          input.Name,
		Type:          input.Type,
		DiscountValue: input.DiscountValue,
		MinAmount:     input.MinAmount,
		MaxDiscount:   input.MaxDiscount,
		TotalCount:    input.TotalCount,
		PerUserLimit:  input.PerUserLimit,
		ValidType:     input.ValidType,
		ValidDays:     input.ValidDays,
		Status:        input.Status,
		IsNewUser:     input.IsNewUser,
		Description:   input.Description,
	}

	if coupon.Status == 0 {
		coupon.Status = 1
	}

	// 解析时间
	if input.StartTime != nil && *input.StartTime != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", *input.StartTime); err == nil {
			localTime := types.LocalTime(t)
			coupon.StartTime = &localTime
		}
	}
	if input.EndTime != nil && *input.EndTime != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", *input.EndTime); err == nil {
			localTime := types.LocalTime(t)
			coupon.EndTime = &localTime
		}
	}

	if err := database.GetDB().Create(&coupon).Error; err != nil {
		utils.Fail(c, "创建优惠券失败")
		return
	}

	utils.Success(c, coupon)
}

// AdminUpdateCoupon 更新优惠券
func AdminUpdateCoupon(c *gin.Context) {
	id := c.Param("id")

	var coupon models.Coupon
	if err := database.GetDB().First(&coupon, id).Error; err != nil {
		utils.Fail(c, "优惠券不存在")
		return
	}

	var input struct {
		Name          string  `json:"name"`
		Type          int8    `json:"type" binding:"min=1,max=3"`
		DiscountValue float64 `json:"discount_value" binding:"min=0"`
		MinAmount     float64 `json:"min_amount"`
		MaxDiscount   float64 `json:"max_discount"`
		TotalCount    int     `json:"total_count"`
		PerUserLimit  int     `json:"per_user_limit" binding:"min=1"`
		ValidType     int8    `json:"valid_type" binding:"min=1,max=2"`
		StartTime     *string `json:"start_time"`
		EndTime       *string `json:"end_time"`
		ValidDays     int     `json:"valid_days"`
		Status        int8    `json:"status"`
		IsNewUser     int8    `json:"is_new_user"`
		Description   string  `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Fail(c, "参数错误: "+err.Error())
		return
	}

	// 更新字段
	if input.Name != "" {
		coupon.Name = input.Name
	}
	if input.Type > 0 {
		coupon.Type = input.Type
	}
	coupon.DiscountValue = input.DiscountValue
	coupon.MinAmount = input.MinAmount
	coupon.MaxDiscount = input.MaxDiscount
	coupon.TotalCount = input.TotalCount
	coupon.PerUserLimit = input.PerUserLimit
	coupon.ValidType = input.ValidType
	coupon.ValidDays = input.ValidDays
	if input.Status >= 0 {
		coupon.Status = input.Status
	}
	coupon.IsNewUser = input.IsNewUser
	coupon.Description = input.Description

	// 解析时间
	if input.StartTime != nil {
		if *input.StartTime == "" {
			coupon.StartTime = nil
		} else if t, err := time.Parse("2006-01-02 15:04:05", *input.StartTime); err == nil {
			localTime := types.LocalTime(t)
			coupon.StartTime = &localTime
		}
	}
	if input.EndTime != nil {
		if *input.EndTime == "" {
			coupon.EndTime = nil
		} else if t, err := time.Parse("2006-01-02 15:04:05", *input.EndTime); err == nil {
			localTime := types.LocalTime(t)
			coupon.EndTime = &localTime
		}
	}

	if err := database.GetDB().Save(&coupon).Error; err != nil {
		utils.Fail(c, "更新优惠券失败")
		return
	}

	utils.Success(c, coupon)
}

// AdminDeleteCoupon 删除优惠券
func AdminDeleteCoupon(c *gin.Context) {
	id := c.Param("id")

	// 检查是否有用户已领取
	var count int64
	database.GetDB().Model(&models.UserCoupon{}).Where("coupon_id = ?", id).Count(&count)
	if count > 0 {
		utils.Fail(c, "该优惠券已有用户领取，无法删除")
		return
	}

	if err := database.GetDB().Delete(&models.Coupon{}, id).Error; err != nil {
		utils.Fail(c, "删除优惠券失败")
		return
	}

	utils.Success(c, gin.H{"message": "删除成功"})
}

// AdminUpdateCouponStatus 更新优惠券状态
func AdminUpdateCouponStatus(c *gin.Context) {
	id := c.Param("id")

	var input struct {
		Status int8 `json:"status" binding:"required,min=0,max=1"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Fail(c, "参数错误")
		return
	}

	if err := database.GetDB().Model(&models.Coupon{}).Where("id = ?", id).Update("status", input.Status).Error; err != nil {
		utils.Fail(c, "更新状态失败")
		return
	}

	utils.Success(c, gin.H{"message": "更新成功"})
}

// AdminGrantCouponToUser 后台给指定用户赠送优惠券
func AdminGrantCouponToUser(c *gin.Context) {
	var input struct {
		CouponID uint64   `json:"coupon_id" binding:"required"`
		UserIDs  []uint64 `json:"user_ids" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Fail(c, "参数错误: "+err.Error())
		return
	}

	// 检查优惠券是否存在
	var coupon models.Coupon
	if err := database.GetDB().First(&coupon, input.CouponID).Error; err != nil {
		utils.Fail(c, "优惠券不存在")
		return
	}

	// 检查优惠券是否启用
	if !coupon.IsAvailable() {
		utils.Fail(c, "优惠券不可用")
		return
	}

	tx := database.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	successCount := 0
	failedUsers := []string{}

	for _, userID := range input.UserIDs {
		// 检查用户是否存在
		var user models.User
		if err := tx.First(&user, userID).Error; err != nil {
			failedUsers = append(failedUsers, strconv.FormatUint(userID, 10)+"(用户不存在)")
			continue
		}

		// 检查用户已领取数量
		var userCount int64
		tx.Model(&models.UserCoupon{}).Where("coupon_id = ? AND user_id = ?", input.CouponID, userID).Count(&userCount)
		if int(userCount) >= coupon.PerUserLimit {
			failedUsers = append(failedUsers, strconv.FormatUint(userID, 10)+"(已达领取上限)")
			continue
		}

		// 计算过期时间
		var expireTime time.Time
		if coupon.ValidType == 1 {
			// 固定时间，使用优惠券的结束时间
			if coupon.EndTime != nil {
				expireTime = time.Time(*coupon.EndTime)
			} else {
				failedUsers = append(failedUsers, strconv.FormatUint(userID, 10)+"(优惠券无有效时间)")
				continue
			}
		} else {
			// 领取后N天
			expireTime = time.Now().AddDate(0, 0, coupon.ValidDays)
		}

		utils.Info("赠送优惠券 - UserID: %d, CouponID: %d, ExpireTime: %s", userID, input.CouponID, expireTime.Format("2006-01-02 15:04:05"))

		// 创建用户优惠券
		userCoupon := models.UserCoupon{
			CouponID:   input.CouponID,
			UserID:     userID,
			Status:     1,
			ExpireTime: types.LocalTime(expireTime),
		}

		if err := tx.Create(&userCoupon).Error; err != nil {
			failedUsers = append(failedUsers, strconv.FormatUint(userID, 10)+"(发放失败)")
			continue
		}

		successCount++
	}

	// 更新优惠券领取数量
	if successCount > 0 {
		if err := tx.Model(&models.Coupon{}).Where("id = ?", input.CouponID).
			Update("received_count", database.GetDB().Raw("received_count + ?", successCount)).Error; err != nil {
			tx.Rollback()
			utils.Fail(c, "更新领取数量失败")
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		utils.Fail(c, "操作失败")
		return
	}

	result := gin.H{
		"success_count": successCount,
		"failed_count":  len(failedUsers),
	}
	if len(failedUsers) > 0 {
		result["failed_users"] = failedUsers
	}

	utils.Success(c, result)
}

// ============ 用户端接口 ============

// GetAvailableCoupons 获取可领取的优惠券列表
func GetAvailableCoupons(c *gin.Context) {
	userID := c.GetUint64("user_id")

	// 查询启用且有效的优惠券
	var coupons []models.Coupon
	query := database.GetDB().Where("status = ?", 1)

	// 如果传了新人券标识，筛选新人券
	isNewUser := c.Query("is_new_user")
	if isNewUser == "1" {
		query = query.Where("is_new_user = ?", 1)
	}

	query.Find(&coupons)

	// 过滤掉用户已领取的
	var availableCoupons []models.Coupon
	for _, coupon := range coupons {
		// 检查优惠券是否可用
		if !coupon.IsAvailable() {
			continue
		}

		// 检查用户是否已领取
		var count int64
		database.GetDB().Model(&models.UserCoupon{}).
			Where("coupon_id = ? AND user_id = ?", coupon.ID, userID).Count(&count)
		if count >= int64(coupon.PerUserLimit) {
			continue
		}

		availableCoupons = append(availableCoupons, coupon)
	}

	utils.Success(c, availableCoupons)
}

// ReceiveCoupon 领取优惠券
func ReceiveCoupon(c *gin.Context) {
	userID := c.GetUint64("user_id")
	couponID := c.Param("id")

	// 检查优惠券是否存在
	var coupon models.Coupon
	if err := database.GetDB().First(&coupon, couponID).Error; err != nil {
		utils.Fail(c, "优惠券不存在")
		return
	}

	// 检查优惠券是否可用
	if !coupon.IsAvailable() {
		utils.Fail(c, "优惠券不可领取")
		return
	}

	// 事务操作（将检查逻辑移入事务内，防止并发问题）
	tx := database.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 在事务内检查用户已领取数量（防止并发重复领取）
	var userCount int64
	tx.Model(&models.UserCoupon{}).
		Where("coupon_id = ? AND user_id = ?", couponID, userID).Count(&userCount)
	if userCount >= int64(coupon.PerUserLimit) {
		tx.Rollback()
		utils.Fail(c, "已达到领取上限")
		return
	}

	// 计算过期时间
	var expireTime time.Time
	if coupon.ValidType == 1 {
		if coupon.EndTime != nil {
			expireTime = time.Time(*coupon.EndTime)
		} else {
			tx.Rollback()
			utils.Fail(c, "优惠券配置错误")
			return
		}
	} else {
		expireTime = time.Now().AddDate(0, 0, coupon.ValidDays)
	}

	// 创建用户优惠券
	userCoupon := models.UserCoupon{
		CouponID:   coupon.ID,
		UserID:     userID,
		Status:     1,
		ExpireTime: types.LocalTime(expireTime),
	}

	if err := tx.Create(&userCoupon).Error; err != nil {
		tx.Rollback()
		// 如果是重复领取（理论上不会发生，因为前面有检查），返回友好提示
		utils.Fail(c, "领取失败，请勿重复操作")
		return
	}

	// 更新优惠券领取数量（使用原子操作防止超发）
	if coupon.TotalCount > 0 {
		result := tx.Model(&models.Coupon{}).
			Where("id = ? AND received_count < total_count", couponID).
			UpdateColumn("received_count", database.GetDB().Raw("received_count + 1"))
		if result.RowsAffected == 0 {
			tx.Rollback()
			utils.Fail(c, "优惠券已领完")
			return
		}
	} else {
		// 不限制总量，直接增加
		tx.Model(&models.Coupon{}).Where("id = ?", couponID).
			UpdateColumn("received_count", database.GetDB().Raw("received_count + 1"))
	}

	if err := tx.Commit().Error; err != nil {
		utils.Fail(c, "领取失败")
		return
	}

	utils.Success(c, gin.H{
		"message":     "领取成功",
		"user_coupon": userCoupon,
	})
}

// GetMyCoupons 获取我的优惠券列表
func GetMyCoupons(c *gin.Context) {
	userID := c.GetUint64("user_id")
	status := c.Query("status") // 1-未使用, 2-已使用, 3-已过期

	query := database.GetDB().Preload("Coupon").Model(&models.UserCoupon{}).
		Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var userCoupons []models.UserCoupon
	if err := query.Order("receive_time DESC").Find(&userCoupons).Error; err != nil {
		utils.Fail(c, "获取优惠券列表失败")
		return
	}

	// 检查并更新过期优惠券
	now := time.Now()
	for i := range userCoupons {
		if userCoupons[i].Status == 1 && now.After(time.Time(userCoupons[i].ExpireTime)) {
			// 标记为已过期
			database.GetDB().Model(&models.UserCoupon{}).
				Where("id = ?", userCoupons[i].ID).
				Update("status", 3)
			userCoupons[i].Status = 3
		}
	}

	utils.Success(c, userCoupons)
}

// GetUsableCoupons 获取可用优惠券（下单时）
func GetUsableCoupons(c *gin.Context) {
	userID := c.GetUint64("user_id")
	orderAmountStr := c.Query("amount")

	if orderAmountStr == "" {
		utils.Fail(c, "请提供订单金额")
		return
	}

	orderAmount, err := strconv.ParseFloat(orderAmountStr, 64)
	if err != nil {
		utils.Fail(c, "订单金额格式错误")
		return
	}

	// 查询用户未使用且未过期的优惠券
	var userCoupons []models.UserCoupon
	database.GetDB().Preload("Coupon").
		Where("user_id = ? AND status = ? AND expire_time > ?", userID, 1, time.Now()).
		Find(&userCoupons)

	// 过滤出满足最低消费金额的优惠券
	var usableCoupons []models.UserCoupon
	for _, uc := range userCoupons {
		if uc.Coupon != nil && orderAmount >= uc.Coupon.MinAmount {
			usableCoupons = append(usableCoupons, uc)
		}
	}

	utils.Success(c, usableCoupons)
}

// CalculateCouponDiscount 计算优惠券优惠金额
func CalculateCouponDiscount(userCouponID uint64, orderAmount float64) (float64, error) {
	var userCoupon models.UserCoupon
	if err := database.GetDB().Preload("Coupon").First(&userCoupon, userCouponID).Error; err != nil {
		return 0, err
	}

	// 检查是否可用
	if !userCoupon.IsUsable() {
		return 0, nil
	}

	coupon := userCoupon.Coupon
	if coupon == nil {
		return 0, nil
	}

	// 检查最低消费
	if orderAmount < coupon.MinAmount {
		return 0, nil
	}

	var discountAmount float64

	switch coupon.Type {
	case 1: // 满减券
		discountAmount = coupon.DiscountValue
	case 2: // 折扣券
		discountAmount = orderAmount * (1 - coupon.DiscountValue)
		if coupon.MaxDiscount > 0 && discountAmount > coupon.MaxDiscount {
			discountAmount = coupon.MaxDiscount
		}
	case 3: // 无门槛券
		discountAmount = coupon.DiscountValue
	}

	// 优惠金额不能超过订单金额
	if discountAmount > orderAmount {
		discountAmount = orderAmount
	}

	return discountAmount, nil
}
