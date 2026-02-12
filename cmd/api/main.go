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
	_ = db.AutoMigrate(&domain.User{})

	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	authService := service.NewAuthService(userRepo, cfg)
	authHandler := handler.NewAuthHandler(authService)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)
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
	{
		protected.GET("/me", authMiddleware.Handler(), authHandler.GetProfile)
		protected.GET("/products", authMiddleware.Handler(), productHandler.GetAllProducts)
		protected.GET("/product/:id", authMiddleware.Handler(), productHandler.GetProductById)
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
			adminProduct.POST("", productHandler.CreateProduct)
			adminProduct.PUT("/:id", productHandler.UpdateProduct)
			adminProduct.PATCH("/:id/deactivate", productHandler.DeactivateProduct)
			adminProduct.PATCH("/:id/activate", productHandler.ActivateProduct)
			adminProduct.DELETE("/:id", productHandler.DeleteProduct)
		}
	}

	superAdminRoutes := router.Group("/api/superadmin")
	superAdminRoutes.Use(authMiddleware.Handler())
	superAdminRoutes.Use(authMiddleware.RequireRoles("super_admin"))
	{
		superAdminRoutes.DELETE("/users/:id/hard", userHandler.HardDeleteUser)
		superAdminRoutes.DELETE("/product/:id/hard", productHandler.HardDeleteProduct)
	}

	log.Printf("Server started on http://localhost:%s", cfg.Port)
	err = router.Run(":" + cfg.Port)
	if err != nil {
		return
	}
}
