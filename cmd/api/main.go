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

	public := r.Group("/")
	{
		public.GET("/products", handlers.GetProducts(db))
		public.GET("/products/campaign", handlers.GetCampaignProducts(db))
		public.GET("/categories", handlers.GetCategories(db))
	}

	auth := r.Group("/auth")
	{
		auth.POST("/register", handlers.Register(db))
		auth.POST("/login", handlers.Login(db, cfg.JWTSecret, cfg.AccessTokenTTL, cfg.RefreshTokenTTL))
		auth.POST("/refresh", handlers.Refresh(db, cfg.JWTSecret, cfg.AccessTokenTTL, cfg.RefreshTokenTTL))
		auth.POST("/logout", handlers.Logout(db))
	}

	r.POST("/admin/login", handlers.AdminLogin(db, cfg.JWTSecret, cfg.AccessTokenTTL))

	admin := r.Group("/admin")
	admin.Use(middleware.AdminAuthMiddleware(cfg.JWTSecret))
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
