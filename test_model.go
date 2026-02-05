package main

import (
	"Inventory-pro/config"
	"Inventory-pro/internal/domain"
	"Inventory-pro/internal/repository"
	"Inventory-pro/pkg/database"
	"fmt"
)

func main() {
	cfg := config.Load()
	db, _ := database.Connect(cfg.DatabaseURL)
	err := db.AutoMigrate(&domain.User{})
	if err != nil {
		return
	}

	productRepo := repository.NewProductRepository(db)

	product := &domain.Product{
		ProductName: "Test Product",
		Unit:        "Test Unit",
		MOQ:         5,
		OM:          5,
		Type:        "Test Type",
		OrderCycle:  "Test OrderCycle",
		AuditCycle:  "Test AuditCycle",
	}

	err = productRepo.Create(product)
	if err != nil {
		fmt.Println("Create failed", err)
		return
	}
	fmt.Printf("Product created with ID: %d\n", product.ID)

	fount, _ := productRepo.FindById(product.ID)
	fmt.Printf("P found with ID: %d\n", fount.ID)

	found, _ := productRepo.FindByProductName(product.ProductName)
	fmt.Printf("Product found with ProductName: %s\n", found.ProductName)

	find, _ := productRepo.FindByProductCode(fount.ProductCode)
	fmt.Printf("Product found with ProductCode: %s\n", find.ProductCode)

	fin, _ := productRepo.FindActiveProduct()
	fmt.Printf("Actived product: %d/n", len(fin))

	products, _ := productRepo.FindAll()
	fmt.Printf("Products found with ID: %d\n", len(products))

}
