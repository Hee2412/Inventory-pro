// test_user_repo.go (root)
package main

import (
	"Inventory-pro/config"
	"Inventory-pro/internal/domain"
	"Inventory-pro/internal/repository"
	"Inventory-pro/pkg/database"
	"Inventory-pro/pkg/password"
	"fmt"
	"log"
)

func main() {
	cfg := config.Load()
	db, _ := database.Connect(cfg.DatabaseURL)

	// Init repository
	userRepo := repository.NewUserRepository(db)

	hashedPassword, err := password.HashPassword("password123")
	if err != nil {
		log.Fatal(err, "Failed to hash password")
	}

	user := domain.User{
		Username:  "user",
		Password:  hashedPassword,
		StoreCode: "ST00001",
	}
	if err := userRepo.Create(&user); err != nil {
		log.Fatal(err, "Failed to create user")
	} else {
		fmt.Println("✅ Created user ID:", user.ID)
	}
}
