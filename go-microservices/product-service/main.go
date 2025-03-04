package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Tạo cấu trúc dữ liệu product
type product struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Remain int    `json:"remain"`
}

// Tạo dữ liệu (sử dụng slice thay vì csdl cho đơn giản)
var products = []product{
	{ID: "1", Name: "Product 1", Remain: 123},
	{ID: "2", Name: "Product 2", Remain: 67},
	{ID: "3", Name: "Product 3", Remain: 10},
}

func main() {
	r := gin.Default()

	// Lấy danh sách product
	r.GET("/products", func(c *gin.Context) {
		c.JSON(http.StatusOK, products)
	})

	// Lấy product theo ID
	r.GET("/products/:id", func(c *gin.Context) {
		id := c.Param("id")
		for _, product := range products {
			if product.ID == id {
				c.JSON(http.StatusOK, product)
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
	})

	// Thêm product mới
	r.POST("/products", func(c *gin.Context) {
		// Lấy dữ liệu từ JSON request
		var newProduct product
		if err := c.ShouldBindJSON(&newProduct); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Kiểm tra trùng ID
		for _, product := range products {
			if product.ID == newProduct.ID {
				c.JSON(http.StatusConflict, gin.H{"error": "Product ID already exists"})
				return
			}
		}

		// Thêm product mới vào danh sách
		products = append(products, newProduct)
		c.JSON(http.StatusCreated, newProduct)
	})

	// Cập nhật product theo ID
	r.PUT("/products/:id", func(c *gin.Context) {
		id := c.Param("id")
		var updatedProduct product

		// Đọc dữ liệu từ request body
		if err := c.ShouldBindJSON(&updatedProduct); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Tìm product cần cập nhật
		for i, product := range products {
			if product.ID == id {
				products[i] = updatedProduct
				c.JSON(http.StatusOK, products[i])
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
	})

	// Xóa product theo ID
	r.DELETE("/products/:id", func(c *gin.Context) {
		id := c.Param("id")

		// Tìm index của product trong danh sách
		for i, product := range products {
			if product.ID == id {
				// Xóa product bằng cách tạo slice mới không chứa product đó
				products = append(products[:i], products[i+1:]...)
				c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
	})

	fmt.Println("Product Service running on port 5003...")
	r.Run(":5003")
}
