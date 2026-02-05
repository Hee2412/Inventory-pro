package main

import (
	"Inventory-pro/config"
	"Inventory-pro/internal/domain"
	"Inventory-pro/pkg/database"
	"fmt"

	"log"
)

func main() {
	// Load config
	cfg := config.Load()

	// Connect database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("❌ Failed to connect:", err)
	}

	fmt.Println("🔄 Migrating all models...")

	// Migrate tất cả models
	err = db.AutoMigrate(
		&domain.User{},
		&domain.Product{},
		&domain.AuditSession{},
		&domain.StoreAuditReport{},
		&domain.OrderSession{},
		&domain.OrderSessionProducts{},
		&domain.StoreOrders{},
		&domain.OrderItems{},
		&domain.SystemLogs{},
		&domain.SystemSettings{},
	)

	if err != nil {
		log.Fatal("❌ Migration failed:", err)
	}

	fmt.Println("✅ All models migrated successfully!")

	// List all tables
	var tables []string
	db.Raw(`
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public' 
		ORDER BY table_name
	`).Scan(&tables)

	fmt.Println("\n Tables created:")
	for i, table := range tables {
		fmt.Printf("  %d. %s\n", i+1, table)
	}

	// Count records in each table
	fmt.Println("\n Record counts:")
	for _, table := range tables {
		var count int64
		db.Table(table).Count(&count)
		fmt.Printf("  - %-30s: %d records\n", table, count)
	}
}
