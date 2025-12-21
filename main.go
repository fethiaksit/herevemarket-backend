package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/handlers"
	"backend/internal/middleware"
)

func main() {
	config.Load()

	client, err := database.Connect(config.AppEnv.MongoURI)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database(config.AppEnv.DBName)

	log.Println("MongoDB connected to:", db.Name())

	r := gin.Default()

	r.GET("/", handlers.Home())

	r.POST("/auth/register", handlers.Register(db))
	r.POST("/auth/login", handlers.Login(
		db,
		config.AppEnv.JWTSecret,
		config.AppEnv.AccessTokenTTL,
		config.AppEnv.RefreshTokenTTL,
	))
	r.POST("/auth/refresh", handlers.Refresh(
		db,
		config.AppEnv.JWTSecret,
		config.AppEnv.AccessTokenTTL,
		config.AppEnv.RefreshTokenTTL,
	))
	r.POST("/auth/logout", handlers.Logout(db))

	r.POST("/admin/login", handlers.AdminLogin(db, config.AppEnv.JWTSecret, config.AppEnv.AccessTokenTTL))

	r.GET("/products", handlers.GetProducts(db))
	r.GET("/categories", handlers.GetCategories(db))
	r.GET("/products/campaign", handlers.GetCampaignProducts(db))

	admin := r.Group("/admin")
	admin.Use(middleware.AdminAuth(config.AppEnv.JWTSecret))
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
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
