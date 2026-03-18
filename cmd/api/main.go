package main

import (
	"Inventory-pro/config"
	"Inventory-pro/internal/domain"
	"Inventory-pro/internal/dto/response"
	"Inventory-pro/internal/handler"
	"Inventory-pro/internal/middleware"
	"Inventory-pro/internal/repository"
	"Inventory-pro/internal/service"
	"Inventory-pro/pkg/database"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect database", err)
	}
	err = db.AutoMigrate(
		&domain.User{}, &domain.Product{},
		&domain.StoreOrder{}, &domain.OrderSession{},
		&domain.OrderSessionProduct{}, &domain.OrderItems{},
		&domain.AuditSession{}, &domain.StoreAuditReport{},
	)
	if err != nil {
		log.Fatal("Failed to auto migrate database", err)
	}
	log.Printf("Database migrated")

	var count int64
	db.Model(&domain.User{}).Where("username = ?", "superadmin").Count(&count)
	if count == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("superadmin123"), bcrypt.DefaultCost)
		superAdmin := domain.User{
			Username: "superadmin",
			Password: string(hashedPassword),
			Role:     "super_admin",
			IsActive: true,
		}
		db.Create(&superAdmin)
		log.Println("✅ Superadmin created!")
	}

	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderSessionRepo := repository.NewOrderSessionRepository(db)
	orderSessionProductRepo := repository.NewOrderSessionProductRepository(db)
	storeOrderRepo := repository.NewStoreOrderRepository(db)
	storeOrderItemRepo := repository.NewStoreOrderItems(db)
	auditSessionRepo := repository.NewAuditSessionRepository(db)
	storeAuditRepo := repository.NewStoreAuditRepository(db)

	authService := service.NewAuthService(userRepo, cfg)
	authHandler := handler.NewAuthHandler(authService)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	orderSessionService := service.NewOrderSessionService(orderSessionRepo, productRepo, orderSessionProductRepo)
	orderSessionHandler := handler.NewOrderSessionHandler(orderSessionService)
	storeOrderService := service.NewStoreOrderService(storeOrderRepo, orderSessionRepo, storeOrderItemRepo, productRepo)
	storeOrderHandler := handler.NewStoreOrderHandler(storeOrderService)
	adminOrderService := service.NewAdminOrderService(orderSessionRepo, storeOrderRepo, userRepo)
	adminOrderHandler := handler.NewAdminOrderHandler(adminOrderService)

	auditSessionService := service.NewAuditSessionService(auditSessionRepo, storeAuditRepo, userRepo, productRepo)
	auditSessionHandler := handler.NewAuditSessionHandler(auditSessionService)
	storeAuditService := service.NewStoreAuditService(auditSessionRepo, storeAuditRepo, userRepo)
	storeAuditHandler := handler.NewStoreAuditHandler(storeAuditService)
	superadminAuditService := service.NewSuperAdminAuditService(auditSessionRepo, storeAuditRepo, userRepo)
	superadminAuditHandler := handler.NewSuperadminAuditHandler(superadminAuditService)

	router := gin.Default()
	router.Use(middleware.LoggingMiddleware())
	router.GET("/health", func(c *gin.Context) {
		response.Message(c, "Service is healthy")
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

	storeProtected := protected.Group("/store")
	storeProtected.Use(authMiddleware.RequireRoles("store"))
	{
		//order routes
		storeProtected.GET("/sessions/:sessionId/order", storeOrderHandler.GetOrCreateOrder)
		storeProtected.PUT("/orders/:orderId/items", storeOrderHandler.UpdateOrder)
		storeProtected.GET("/orders/:orderId", storeOrderHandler.GetOrderDetail)
		storeProtected.GET("/orders", storeOrderHandler.GetMyOrder)
		storeProtected.PUT("/orders/:orderId", storeOrderHandler.UpdateStatus)
		//audit routes
		storeProtected.PUT("/audit-sessions/:sessionId/items", storeAuditHandler.UpdateAuditItem)
		storeProtected.GET("/audit-reports", storeAuditHandler.GetMyAuditReport)
		storeProtected.GET("/audit-sessions/:sessionId/report", storeAuditHandler.GetAuditReport)
	}

	adminRoutes := protected.Group("/admin")
	adminRoutes.Use(authMiddleware.RequireRoles("admin", "super_admin"))
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
		adminAudit := adminRoutes.Group("/audit-sessions")
		{
			adminAudit.POST("", auditSessionHandler.CreateAuditSession)
			adminAudit.GET("", auditSessionHandler.GetAllAuditSession)
			adminAudit.GET("/:sessionId", auditSessionHandler.GetAuditSessionByID)
			adminAudit.POST("/products", auditSessionHandler.AddProductToAudit)
			adminAudit.DELETE("/:sessionId/products/:productId", auditSessionHandler.RemoveProductFromAudit)
			adminAudit.PATCH("/:id/close", auditSessionHandler.CloseAuditSession)
			adminAudit.PUT("/:id", auditSessionHandler.UpdateAuditSession)
		}
		adminRoutes.GET("/orders", adminOrderHandler.GetAllOrders)
		adminRoutes.POST("/orders/:orderId/approve", adminOrderHandler.ApproveOrder)
		adminRoutes.POST("/orders/:orderId/decline", adminOrderHandler.DeclineOrder)
	}

	superAdminRoutes := protected.Group("/superadmin")
	superAdminRoutes.Use(authMiddleware.RequireRoles("super_admin"))
	{
		superAdminRoutes.DELETE("/users/:id/hard", userHandler.HardDeleteUser)
		superAdminRoutes.DELETE("/products/:id/hard", productHandler.HardDeleteProduct)

		superAudit := superAdminRoutes.Group("/audit-sessions")
		{
			superAudit.GET("/:sessionId/reports", superadminAuditHandler.GetAllReportsInSession)
			superAudit.GET("/:sessionId/stores/:storeId", superadminAuditHandler.GetReportDetail)
			superAudit.GET("/:sessionId/summary", superadminAuditHandler.GetAuditSummary)
			superAudit.POST("/:sessionId/stores/:storeId/approve", superadminAuditHandler.ApproveStoreReport)
			superAudit.POST("/:sessionId/stores/:storeId/decline", superadminAuditHandler.DeclineStoreReport)
		}
	}

	log.Printf("Server started on http://localhost:%s", cfg.Port)
	err = router.Run(":" + cfg.Port)
	if err != nil {
		return
	}
}
