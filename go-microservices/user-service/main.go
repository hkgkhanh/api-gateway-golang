package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Tạo cấu trúc dữ liệu user
type user struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// Tạo dữ liệu (sử dụng slice thay vì csdl cho đơn giản)
var users = []user{
	{ID: "1", Name: "Nong Quoc Khanh", Age: 20},
	{ID: "2", Name: "Nguyen Van A", Age: 25},
	{ID: "3", Name: "Tran Thi B", Age: 22},
}

func main() {
	r := gin.Default()

	// Lấy danh sách user
	r.GET("/users", func(c *gin.Context) {
		c.JSON(http.StatusOK, users)
	})

	// Lấy user theo ID
	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		for _, user := range users {
			if user.ID == id {
				c.JSON(http.StatusOK, user)
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	})

	// Thêm user mới
	r.POST("/users", func(c *gin.Context) {
		// Lấy dữ liệu từ JSON request
		var newUser user
		if err := c.ShouldBindJSON(&newUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Kiểm tra trùng ID
		for _, u := range users {
			if u.ID == newUser.ID {
				c.JSON(http.StatusConflict, gin.H{"error": "User ID already exists"})
				return
			}
		}

		// Thêm user mới vào danh sách
		users = append(users, newUser)
		c.JSON(http.StatusCreated, newUser)
	})

	// Cập nhật user theo ID
	r.PUT("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		var updatedUser user

		// Đọc dữ liệu từ request body
		if err := c.ShouldBindJSON(&updatedUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Tìm user cần cập nhật
		for i, u := range users {
			if u.ID == id {
				users[i] = updatedUser
				c.JSON(http.StatusOK, users[i])
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	})

	// Xóa user theo ID
	r.DELETE("/users/:id", func(c *gin.Context) {
		id := c.Param("id")

		// Tìm index của user trong danh sách
		for i, u := range users {
			if u.ID == id {
				// Xóa user bằng cách tạo slice mới không chứa user đó
				users = append(users[:i], users[i+1:]...)
				c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	})

	fmt.Println("User Service running on port 5001...")
	r.Run(":5001")
}
