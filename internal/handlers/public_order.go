package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"backend/internal/models"
)

type createOrderItemRequest struct {
	ProductID string  `json:"productId" binding:"required"`
	Name      string  `json:"name" binding:"required"`
	Price     float64 `json:"price" binding:"required"`
	Quantity  int     `json:"quantity" binding:"required"`
}

type createOrderCustomerRequest struct {
	Title  string `json:"title" binding:"required"`
	Detail string `json:"detail" binding:"required"`
	Note   string `json:"note"`
}

type createOrderRequest struct {
	Items         []createOrderItemRequest   `json:"items" binding:"required"`
	TotalPrice    float64                    `json:"totalPrice"`
	Customer      createOrderCustomerRequest `json:"customer" binding:"required"`
	PaymentMethod string                     `json:"paymentMethod" binding:"required"`
}

func CreateOrder(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		const route = "POST /orders"
		defer handlePanic(c, route)

		if err := ensureDBConnection(c.Request.Context(), db); err != nil {
			respondWithError(c, http.StatusServiceUnavailable, route, "database unavailable")
			return
		}

		var req createOrderRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			respondWithError(c, http.StatusBadRequest, route, "invalid request body")
			return
		}

		order, err := buildOrderFromRequest(req)
		if err != nil {
			respondWithError(c, http.StatusBadRequest, route, err.Error())
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		res, err := db.Collection("orders").InsertOne(ctx, order)
		if err != nil {
			respondWithError(c, http.StatusInternalServerError, route, "db error")
			return
		}

		id, ok := res.InsertedID.(primitive.ObjectID)
		if ok {
			order.ID = id
		}

		c.JSON(http.StatusCreated, gin.H{
			"orderId": order.ID.Hex(),
			"message": "order created",
		})
	}
}

func buildOrderFromRequest(req createOrderRequest) (models.Order, error) {
	if len(req.Items) == 0 {
		return models.Order{}, errors.New("at least one item is required")
	}

	if req.PaymentMethod != "cash" && req.PaymentMethod != "card" {
		return models.Order{}, errors.New("invalid payment method")
	}

	items := make([]models.OrderItem, 0, len(req.Items))
	var total float64

	for _, item := range req.Items {
		productID, err := primitive.ObjectIDFromHex(item.ProductID)
		if err != nil {
			return models.Order{}, errors.New("invalid productId")
		}

		if item.Quantity <= 0 {
			return models.Order{}, errors.New("quantity must be greater than zero")
		}

		if item.Price < 0 {
			return models.Order{}, errors.New("price must be zero or greater")
		}

		items = append(items, models.OrderItem{
			ProductID: productID,
			Name:      item.Name,
			Price:     item.Price,
			Quantity:  item.Quantity,
		})

		total += item.Price * float64(item.Quantity)
	}

	order := models.Order{
		Items:         items,
		TotalPrice:    total,
		Customer:      models.OrderCustomer(req.Customer),
		PaymentMethod: req.PaymentMethod,
		Status:        "pending",
		CreatedAt:     time.Now(),
	}

	return order, nil
}
