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
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":   "OK",
			"database": "Connected",
		})
	})

	mw := middleware.NewAuthMiddleware(cfg)
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/register", mw.Handler(), authHandler.Register)
	}
	protected := router.Group("/api")
	{
		protected.GET("/me", mw.Handler(), authHandler.GetProfile)
	}

	log.Printf("Server started on http://localhost:%s", cfg.Port)
	err = router.Run(":" + cfg.Port)
	if err != nil {
		return
	}
}
