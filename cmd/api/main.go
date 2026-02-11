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
	authService := service.NewAuthService(userRepo, cfg)
	authHandler := handler.NewAuthHandler(authService)
	userService := service.NewUserService(userRepo, cfg)
	userHandler := handler.NewUserHandler(userService)
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
	}

	adminRoutes := router.Group("/api/admin")
	adminRoutes.Use(authMiddleware.Handler())
	adminRoutes.Use(authMiddleware.RequireRoles("admin", "super_admin"))
	{
		adminRoutes.PUT("/register", authHandler.Register)
		adminRoutes.GET("/users", userHandler.GetAllUsers)
		adminRoutes.GET("/users/:id", userHandler.GetUserById)
		adminRoutes.PUT("/users/:id", userHandler.UpdateUser)
		adminRoutes.PATCH("/users/:id/deactivate", userHandler.DeactivateUser)
		adminRoutes.PATCH("/users/:id/activate", userHandler.ActivateUser)
	}
	superAdminRoutes := router.Group("/api/superadmin")
	superAdminRoutes.Use(authMiddleware.Handler())
	superAdminRoutes.Use(authMiddleware.RequireRoles("super_admin"))
	{
		superAdminRoutes.DELETE("/users/:id", userHandler.DeleteUser)
		superAdminRoutes.DELETE("/users/:id/hard", userHandler.HardDeleteUser)
	}

	log.Printf("Server started on http://localhost:%s", cfg.Port)
	err = router.Run(":" + cfg.Port)
	if err != nil {
		return
	}
}
