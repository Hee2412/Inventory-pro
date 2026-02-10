package main

import (
	"Inventory-pro/config"
	"Inventory-pro/internal/domain"
	"Inventory-pro/pkg/database"
	"Inventory-pro/pkg/password"
	"fmt"
)

func main() {
	cfg := config.Load()
	db, _ := database.Connect(cfg.DatabaseURL)

	// Hash password
	hashedPassword, _ := password.HashPassword("superadmin123")

	// Create admin
	admin := &domain.User{
		Username: "superadmin",
		Password: hashedPassword,
		Role:     "super_admin",
		IsActive: true,
	}

	err := db.Create(admin).Error
	if err != nil {
		fmt.Println("❌ Failed:", err)
		return
	}

	fmt.Printf("✅ Admin created! ID: %d\n", admin.ID)
	fmt.Println("Username: admin")
	fmt.Println("Password: admin123")
}
