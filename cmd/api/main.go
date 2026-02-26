package main

import (
	"Inventory-pro/config"
	"Inventory-pro/internal/domain"
	"Inventory-pro/internal/handler"
	"Inventory-pro/internal/middleware"
	"Inventory-pro/internal/repository"
	"Inventory-pro/internal/service"
	"Inventory-pro/pkg/database"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect database", err)
	}
	_ = db.AutoMigrate(
		&domain.User{}, &domain.Product{},
		&domain.StoreOrder{}, &domain.OrderSession{},
		&domain.OrderSessionProduct{}, &domain.OrderItems{},
	)

	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderSessionRepo := repository.NewOrderSessionRepository(db)
	orderSessionProductRepo := repository.NewOrderSessionProductRepository(db)
	storeOrderRepo := repository.NewStoreOrderRepository(db)
	storeOrderItemRepo := repository.NewStoreOrderItems(db)

	authService := service.NewAuthService(userRepo, cfg)
	authHandler := handler.NewAuthHandler(authService)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	orderSessionService := service.NewOrderSessionService(orderSessionRepo, productRepo, orderSessionProductRepo)
	orderSessionHandler := handler.NewOrderSessionHandler(orderSessionService)
	storeOrderService := service.NewStoreOrderService(storeOrderRepo, orderSessionRepo, storeOrderItemRepo)
	storeOrderHandler := handler.NewStoreOrderHandler(storeOrderService)
	adminOrderService := service.NewAdminOrderService(orderSessionRepo, storeOrderRepo)
	adminOrderHandler := handler.NewAdminOrderHandler(adminOrderService)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":   "OK",
			"database": "Connected",
		})
	})

	authMiddleware := middleware.NewAuthMiddleware(cfg)

	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/login", authHandler.Login)
	}

	protected := router.Group("/api")
	protected.Use(authMiddleware.Handler())
	{
		protected.GET("/me", authHandler.GetProfile)
		protected.GET("/products", productHandler.GetAllProducts)
		protected.GET("/products/:id", productHandler.GetProductById)
	}

	storeProtected := protected.Group("api/store")
	storeProtected.Use(authMiddleware.Handler(), authMiddleware.RequireRoles("store"))
	{
		storeProtected.GET("/sessions/:sessionId/order", storeOrderHandler.GetOrCreateOrder)
		storeProtected.PUT("/orders/:orderId/items", storeOrderHandler.UpdateOrder)
		storeProtected.POST("/orders/:orderId/submit", storeOrderHandler.SubmitOrder)
		storeProtected.GET("/orders/:orderId", storeOrderHandler.GetOrderDetail)
		storeProtected.GET("/orders", storeOrderHandler.GetMyOrder)
	}

	adminRoutes := router.Group("/api/admin")
	adminRoutes.Use(authMiddleware.Handler(), authMiddleware.RequireRoles("admin", "super_admin"))
	{
		adminUser := adminRoutes.Group("/users")
		{
			adminUser.POST("/register", authHandler.Register)
			adminUser.GET("", userHandler.GetAllUsers)
			adminUser.GET("/:id", userHandler.GetUserById)
			adminUser.PUT("/:id", userHandler.UpdateUser)
			adminUser.PATCH("/:id/deactivate", userHandler.DeactivateUser)
			adminUser.PATCH("/:id/activate", userHandler.ActivateUser)
			adminUser.DELETE("/:id", userHandler.DeleteUser)
		}

		adminProduct := adminRoutes.Group("/products")
		{
			adminProduct.GET("", productHandler.GetAllProductsForAdmin)
			adminProduct.POST("", productHandler.CreateProduct)
			adminProduct.PUT("/:id", productHandler.UpdateProduct)
			adminProduct.PATCH("/:id/deactivate", productHandler.DeactivateProduct)
			adminProduct.PATCH("/:id/activate", productHandler.ActivateProduct)
			adminProduct.DELETE("/:id", productHandler.DeleteProduct)
		}
		adminSession := adminRoutes.Group("/sessions")
		{
			adminSession.GET("/:sessionId/orders", adminOrderHandler.GetAllOrderInSession)
			adminSession.POST("", orderSessionHandler.CreateSession)
			adminSession.GET("", orderSessionHandler.GetAllSessions)
			adminSession.GET("/:sessionId", orderSessionHandler.GetSessionById)
			adminSession.POST("/products", orderSessionHandler.AddProductToSession)
			adminSession.DELETE("/:sessionId/products/:productId", orderSessionHandler.RemoveProductFromSession)
			adminSession.PATCH("/:sessionId/close", orderSessionHandler.CloseSession)
		}
		adminRoutes.POST("/orders/:orderId/approve", adminOrderHandler.ApproveOrder)
		adminRoutes.POST("/orders/:orderId/decline", adminOrderHandler.DeclineOrder)
	}

	superAdminRoutes := router.Group("/api/superadmin")
	superAdminRoutes.Use(authMiddleware.Handler())
	superAdminRoutes.Use(authMiddleware.RequireRoles("super_admin"))
	{
		superAdminRoutes.DELETE("/users/:id/hard", userHandler.HardDeleteUser)
		superAdminRoutes.DELETE("/products/:id/hard", productHandler.HardDeleteProduct)
	}

	log.Printf("Server started on http://localhost:%s", cfg.Port)
	err = router.Run(":" + cfg.Port)
	if err != nil {
		return
	}
}
