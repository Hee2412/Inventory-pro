package main

import (
	"Inventory-pro/config"
	"Inventory-pro/pkg/database"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	cfg := config.Load()
	log.Printf("Starting server in %s mode", cfg.Environment)

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("False to connect database", err)
	}
	log.Printf("Connected to database successfully")

	var result int
	db.Raw("SELECT 1").Scan(&result)
	log.Printf("Result: %d", result)

	router := gin.Default()

	router.GET("/heath", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":   "OK",
			"database": "Connected",
		})
	})

	log.Printf("Server started on http://localhost:%s", cfg.Port)
	router.Run(":" + cfg.Port)
}
