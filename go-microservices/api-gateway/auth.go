package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Khóa bí mật để ký JWT
var jwtSecret = []byte("secret-key-123")

// Cấu trúc User giả lập
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Danh sách user giả lập (thay bằng database trong thực tế)
var users = map[string]string{
	"admin": "admin123",
	"user":  "user123",
}

var roles = map[string]string{
	"admin": "admin",
	"user":  "user",
}

// Hàm tạo JWT
func generateJWT(username string, role string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * 1).Unix(), // Hết hạn sau 1 giờ
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// API Login để cấp JWT
func loginHandler(c *gin.Context) {
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Kiểm tra username/password
	password, exists := users[loginData.Username]
	if !exists || password != loginData.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Lấy role của user
	role := roles[loginData.Username]

	// Tạo JWT
	token, err := generateJWT(loginData.Username, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
