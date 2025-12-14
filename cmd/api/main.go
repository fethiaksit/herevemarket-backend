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

	r.POST("/admin/login", handlers.AdminLogin(db, cfg.JWTSecret))

	admin := r.Group("/admin")
	admin.Use(middleware.AdminAuth(cfg.JWTSecret))
	{
		admin.GET("/me", func(c *gin.Context) {
			c.JSON(200, gin.H{"ok": true})
		})
		r.GET("/products", handlers.GetProducts(db))
	}

	r.Run(":8080")
}
