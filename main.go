package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"path"

	"go-products.com/m/internal"
	"go-products.com/m/internal/product/infrastructure/persistance"
	"go-products.com/m/internal/product/infrastructure/persistance/migrations"
	"go-products.com/m/internal/shared/database"
)

func main() {
	db, err := database.GenerateDatabaseConnection(database.DatabaseConnection{DatabaseName: "products.db"}, migrations.CreateProductsDatabase)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	productRepository := persistance.NewProductsSQLiteRepository(db)
	err = migrations.InitProducts(context.Background(), productRepository, path.Join(dir, "infra", "migrations", "data.json"))
	if err != nil {
		log.Fatal(err)
	}

	server := internal.SetupServer(productRepository)

	log.Println("Server running on port 8080")
	err = http.ListenAndServe(":8080", server)
	if err != nil {
		log.Fatal(err)
	}
}
