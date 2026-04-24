package routes

import (
	"shop_api/handlers"
	"shop_api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORSMiddleware())

	api := r.Group("/api")
	{
		// 上传接口（需要认证）
		upload := api.Group("/upload")
		upload.Use(middleware.Auth())
		{
			upload.POST("", handlers.UploadImage)       // 上传图片
			upload.POST("/video", handlers.UploadVideo) // 上传视频
		}

		miniapp := api.Group("/miniapp")
		{
			miniapp.POST("/register", handlers.Register)
			miniapp.POST("/login", handlers.Login)
			miniapp.POST("/wechat_login", handlers.WechatLogin)

			miniapp.GET("/banners", handlers.GetBanners)
			miniapp.GET("/categories", handlers.GetCategories)
			miniapp.GET("/categories/:id/sub", handlers.GetSubCategories)
			miniapp.GET("/products", handlers.GetProducts)
			miniapp.GET("/products/recommend", handlers.GetRecommendProducts)
			miniapp.GET("/products/new", handlers.GetNewProducts)
			miniapp.GET("/products/:id", handlers.GetProduct)

			// 系统配置（小程序端）
			miniapp.GET("/config/:name", handlers.GetConfigByName)
			miniapp.POST("/configs/batch", handlers.GetConfigsByNames)

			// 帮助中心（小程序端）
			miniapp.GET("/help/categories", handlers.GetActiveHelpCategories)
			miniapp.GET("/help/questions/:category_id", handlers.GetHelpQuestionsByCategory)
			miniapp.GET("/help/question/:id", handlers.GetHelpQuestionDetail)
			miniapp.GET("/help/search", handlers.SearchHelpQuestions)
			miniapp.GET("/help/config", handlers.GetHelpConfig)

			// 运费计算（小程序端）
			miniapp.POST("/shipping/calculate", handlers.CalculateShippingFee)
		}

		user := api.Group("/user")
		user.Use(middleware.Auth())
		user.Use(middleware.OperationLog()) // 添加操作日志中间件
		{
			user.GET("/info", handlers.GetUserInfo)
			user.PUT("/info", handlers.UpdateUserInfo)
			user.PUT("/password", handlers.ChangePassword)

			user.GET("/addresses", handlers.GetAddresses)
			user.GET("/addresses/:id", handlers.GetAddress)
			user.POST("/addresses", handlers.CreateAddress)
			user.PUT("/addresses/:id", handlers.UpdateAddress)
			user.DELETE("/addresses/:id", handlers.DeleteAddress)
			user.PUT("/addresses/:id/default", handlers.SetDefaultAddress)

			user.GET("/cart", handlers.GetCart)
			user.GET("/cart/count", handlers.GetCartCount)
			user.POST("/cart", handlers.AddCart)
			user.PUT("/cart/:id", handlers.UpdateCart)
			user.PUT("/cart/:id/select", handlers.SelectCart)
			user.PUT("/cart/select_all", handlers.SelectAllCart)
			user.DELETE("/cart/:id", handlers.DeleteCart)
			user.DELETE("/cart", handlers.ClearCart)

			user.GET("/orders", handlers.GetOrders)
			user.GET("/orders/:id", handlers.GetOrder)
			user.POST("/orders", handlers.CreateOrder)
			user.PUT("/orders/:id/cancel", handlers.CancelOrder)
			user.PUT("/orders/:id/confirm", handlers.ConfirmReceive)
			user.DELETE("/orders/:id", handlers.DeleteOrder)

			user.POST("/pay", handlers.GetPayURL)
			user.GET("/pay/status/:id", handlers.QueryPayStatus)
			user.POST("/pay/mock_success/:id", handlers.MockPaySuccess)
			user.POST("/pay/refund/:id", handlers.ApplyRefund)

			// 退款管理（用户端）
			user.GET("/refunds", handlers.GetMyRefunds)
			user.GET("/refunds/:id", handlers.GetRefundDetail)
			user.POST("/refunds/apply", handlers.ApplyRefund)

			// 订单物流查询
			user.GET("/orders/:id/logistics", handlers.GetOrderLogistics)

			// 优惠券（用户端）
			user.GET("/coupons/available", handlers.GetAvailableCoupons)
			user.POST("/coupons/receive/:id", handlers.ReceiveCoupon)
			user.GET("/coupons/my", handlers.GetMyCoupons)
			user.GET("/coupons/usable", handlers.GetUsableCoupons)
		}

		admin := api.Group("/admin")
		{
			admin.POST("/register", handlers.AdminRegister)
			admin.POST("/login", handlers.AdminLogin)
		}

		adminAuth := api.Group("/admin")
		adminAuth.Use(middleware.AdminAuth())
		adminAuth.Use(middleware.OperationLog()) // 添加操作日志中间件
		{
			// 用户管理
			adminAuth.GET("/users", handlers.GetUsers)
			adminAuth.GET("/users/:id", handlers.GetUser)
			adminAuth.PUT("/users/:id", handlers.UpdateUserInfo)
			adminAuth.PUT("/users/:id/status", handlers.UpdateUserStatus)
			adminAuth.DELETE("/users/:id", handlers.DeleteUser)

			// 分类管理
			adminAuth.GET("/categories", handlers.GetCategories)
			adminAuth.GET("/categories/:id", handlers.GetCategory)
			adminAuth.POST("/categories", handlers.CreateCategory)
			adminAuth.PUT("/categories/:id", handlers.UpdateCategory)
			adminAuth.DELETE("/categories/:id", handlers.DeleteCategory)

			// 商品管理
			adminAuth.GET("/products", handlers.GetProducts)
			adminAuth.GET("/products/:id", handlers.GetProduct)
			adminAuth.POST("/products", handlers.CreateProduct)
			adminAuth.PUT("/products/:id", handlers.UpdateProduct)
			adminAuth.DELETE("/products/:id", handlers.DeleteProduct)

			// 轮播图管理
			adminAuth.GET("/banners", handlers.GetBanners)
			adminAuth.POST("/banners", handlers.CreateBanner)
			adminAuth.PUT("/banners/:id", handlers.UpdateBanner)
			adminAuth.DELETE("/banners/:id", handlers.DeleteBanner)

			// 订单管理
			adminAuth.GET("/orders", handlers.AdminGetOrders)
			adminAuth.GET("/orders/:id", handlers.AdminGetOrder)
			adminAuth.PUT("/orders/:id/ship", handlers.ShipOrder)
			adminAuth.PUT("/orders/:id/status", handlers.UpdateOrderStatus)

			// 系统配置管理
			adminAuth.GET("/configs", handlers.GetConfigs)
			adminAuth.GET("/configs/:id", handlers.GetConfig)
			adminAuth.POST("/configs", handlers.CreateConfig)
			adminAuth.PUT("/configs/:id", handlers.UpdateConfig)
			adminAuth.DELETE("/configs/:id", handlers.DeleteConfig)

			// 帮助中心分类管理
			adminAuth.GET("/help/categories", handlers.GetHelpCategories)
			adminAuth.POST("/help/categories", handlers.CreateHelpCategory)
			adminAuth.PUT("/help/categories/:id", handlers.UpdateHelpCategory)
			adminAuth.DELETE("/help/categories/:id", handlers.DeleteHelpCategory)

			// 帮助中心问题管理
			adminAuth.GET("/help/questions", handlers.GetHelpQuestions)
			adminAuth.GET("/help/questions/:id", handlers.GetHelpQuestionDetail)
			adminAuth.POST("/help/questions", handlers.CreateHelpQuestion)
			adminAuth.PUT("/help/questions/:id", handlers.UpdateHelpQuestion)
			adminAuth.DELETE("/help/questions/:id", handlers.DeleteHelpQuestion)
			adminAuth.PUT("/help/questions/batch-status", handlers.BatchUpdateQuestionsStatus)
			adminAuth.GET("/help/statistics", handlers.GetHelpStatistics)

			// 优惠券管理（后台）
			adminAuth.GET("/coupons", handlers.AdminCouponList)
			adminAuth.GET("/coupons/:id", handlers.AdminCouponDetail)
			adminAuth.POST("/coupons", handlers.AdminCreateCoupon)
			adminAuth.PUT("/coupons/:id", handlers.AdminUpdateCoupon)
			adminAuth.DELETE("/coupons/:id", handlers.AdminDeleteCoupon)
			adminAuth.POST("/coupons/:id/status", handlers.AdminUpdateCouponStatus)
			adminAuth.POST("/coupons/grant", handlers.AdminGrantCouponToUser)

			// 运费模板管理（后台）
			adminAuth.GET("/shipping/templates", handlers.GetShippingTemplates)
			adminAuth.GET("/shipping/templates/:id", handlers.GetShippingTemplate)
			adminAuth.POST("/shipping/templates", handlers.CreateShippingTemplate)
			adminAuth.PUT("/shipping/templates/:id", handlers.UpdateShippingTemplate)
			adminAuth.DELETE("/shipping/templates/:id", handlers.DeleteShippingTemplate)
			adminAuth.POST("/shipping/templates/:id/default", handlers.SetDefaultTemplate)

			// 退款管理（后台）
			adminAuth.GET("/refunds", handlers.AdminGetRefunds)
			adminAuth.GET("/refunds/:id", handlers.AdminGetRefundDetail)
			adminAuth.POST("/refunds/:id/approve", handlers.AdminApproveRefund)
			adminAuth.POST("/refunds/:id/reject", handlers.AdminRejectRefund)

			// 数据统计（后台）
			adminAuth.GET("/statistics/dashboard", handlers.GetDashboardStats)
			adminAuth.GET("/statistics/sales-trend", handlers.GetSalesTrend)
			adminAuth.GET("/statistics/users-trend", handlers.GetUsersTrend)
		}
	}

	r.POST("/api/notify/wechat", handlers.WechatNotify)
	r.POST("/api/notify/alipay", handlers.AlipayNotify)

	r.Static("/uploads", "./uploads")

	return r
}
