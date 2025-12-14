package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/handlers"
	"backend/internal/middleware"
)

func main() {
	cfg := config.Load()

	log.Println("MONGO_URI =", cfg.MongoURI)

	client, err := database.Connect(cfg.MongoURI)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database(cfg.DBName)

	r := gin.Default()

	r.GET("/", handlers.Home())

	r.POST("/admin/login", handlers.AdminLogin(db, cfg.JWTSecret))

	r.GET("/products", handlers.GetProducts(db))
	r.GET("/categories", handlers.GetCategories(db))

	admin := r.Group("/admin")
	admin.Use(middleware.AdminAuth(cfg.JWTSecret))
	{
		admin.GET("/me", func(c *gin.Context) {
			c.JSON(200, gin.H{"ok": true})
		})

		admin.GET("/products", handlers.GetAllProducts(db))
		admin.POST("/products", handlers.CreateProduct(db))
		admin.PUT("/products/:id", handlers.UpdateProduct(db))
		admin.DELETE("/products/:id", handlers.DeleteProduct(db))

		admin.GET("/categories", handlers.GetAllCategories(db))
		admin.POST("/categories", handlers.CreateCategory(db))
		admin.PUT("/categories/:id", handlers.UpdateCategory(db))
		admin.DELETE("/categories/:id", handlers.DeleteCategory(db))
	}

	r.Run(":8080")
}
