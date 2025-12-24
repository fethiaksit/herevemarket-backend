package handlers

import "github.com/gin-gonic/gin"

func AdminLoginPage(c *gin.Context) {
	c.HTML(200, "admin/login.html", gin.H{})
}

func AdminCategoriesPage(c *gin.Context) {
	c.HTML(200, "admin/categories.html", gin.H{})
}

func AdminProductsPage(c *gin.Context) {
	c.HTML(200, "admin/products.html", gin.H{})
}

func AdminOrdersPage(c *gin.Context) {
	c.HTML(200, "admin/orders.html", gin.H{})
}
