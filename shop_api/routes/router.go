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
			adminAuth.PUT("/users/:id/status", handlers.UpdateUserStatus)
			adminAuth.DELETE("/users/:id", handlers.DeleteUser)

			adminAuth.GET("/categories", handlers.GetCategories)
			adminAuth.GET("/categories/:id", handlers.GetCategory)
			adminAuth.POST("/categories", handlers.CreateCategory)
			adminAuth.PUT("/categories/:id", handlers.UpdateCategory)
			adminAuth.DELETE("/categories/:id", handlers.DeleteCategory)

			adminAuth.GET("/products", handlers.GetProducts)
			adminAuth.GET("/products/:id", handlers.GetProduct)
			adminAuth.POST("/products", handlers.CreateProduct)
			adminAuth.PUT("/products/:id", handlers.UpdateProduct)
			adminAuth.DELETE("/products/:id", handlers.DeleteProduct)

			adminAuth.GET("/banners", handlers.GetBanners)
			adminAuth.POST("/banners", handlers.CreateBanner)
			adminAuth.PUT("/banners/:id", handlers.UpdateBanner)
			adminAuth.DELETE("/banners/:id", handlers.DeleteBanner)
		}
	}

	r.POST("/api/notify/wechat", handlers.WechatNotify)
	r.POST("/api/notify/alipay", handlers.AlipayNotify)

	r.Static("/uploads", "./uploads")

	return r
}
