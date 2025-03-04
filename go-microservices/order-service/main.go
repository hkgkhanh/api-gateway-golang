package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Tạo cấu trúc dữ liệu order
type order struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	ProductID string `json:"productid"`
	Quantity  int    `json:"quantity"`
}

// Tạo dữ liệu slice thay vì csdl cho đơn giản
var orders = []order{
	{ID: "1", Title: "Order ABC", ProductID: "1", Quantity: 2},
	{ID: "2", Title: "Order MNP", ProductID: "3", Quantity: 5},
	{ID: "3", Title: "Order XYZ", ProductID: "2", Quantity: 1},
}

func main() {
	r := gin.Default()

	// Lấy danh sách order
	r.GET("/orders", func(c *gin.Context) {
		c.JSON(http.StatusOK, orders)
	})

	// Lấy order theo ID
	r.GET("/orders/:id", func(c *gin.Context) {
		id := c.Param("id")
		for _, order := range orders {
			if order.ID == id {
				c.JSON(http.StatusOK, order)
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
	})

	// Thêm order mới
	r.POST("/orders", func(c *gin.Context) {
		// Lấy dữ liệu từ JSON request
		var newOrder order
		if err := c.ShouldBindJSON(&newOrder); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Kiểm tra trùng ID
		for _, order := range orders {
			if order.ID == newOrder.ID {
				c.JSON(http.StatusConflict, gin.H{"error": "Order ID already exists"})
				return
			}
		}

		// Thêm order mới vào danh sách
		orders = append(orders, newOrder)
		c.JSON(http.StatusCreated, newOrder)
	})

	// Cập nhật order theo ID
	r.PUT("/orders/:id", func(c *gin.Context) {
		id := c.Param("id")
		var updatedOrder order

		// Đọc dữ liệu từ request body
		if err := c.ShouldBindJSON(&updatedOrder); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Tìm order cần cập nhật
		for i, order := range orders {
			if order.ID == id {
				orders[i] = updatedOrder
				c.JSON(http.StatusOK, orders[i])
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
	})

	// Xóa order theo ID
	r.DELETE("/orders/:id", func(c *gin.Context) {
		id := c.Param("id")

		// Tìm index của order trong danh sách
		for i, order := range orders {
			if order.ID == id {
				// Xóa orders bằng cách tạo slice mới không chứa order đó
				orders = append(orders[:i], orders[i+1:]...)
				c.JSON(http.StatusOK, gin.H{"message": "Order deleted"})
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
	})

	fmt.Println("Order Service running on port 5002...")
	r.Run(":5002")
}
