package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// // Middleware kiểm tra JWT
// func authMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		authHeader := c.GetHeader("Authorization")
// 		if authHeader == "" {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization Header"})
// 			c.Abort()
// 			return
// 		}

// 		// Token có dạng "Bearer <token>"
// 		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
// 		if tokenString == authHeader {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization Header"})
// 			c.Abort()
// 			return
// 		}

// 		// Xác thực token
// 		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 			return jwtSecret, nil
// 		})

// 		if err != nil || !token.Valid {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
// 			c.Abort()
// 			return
// 		}

// 		// Nếu hợp lệ, tiếp tục request
// 		c.Next()
// 	}
// }

// // Giả lập danh sách token hợp lệ
// var validTokens = map[string]string{
// 	"admin-token": "admin",
// 	"user-token":  "user",
// }

// Middleware xác thực JWT
func authMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		// Loại bỏ prefix "Bearer "
		tokenString = tokenString[7:]

		// Giải mã JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Kiểm tra role (Authorization)
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		userRole := claims["role"].(string)
		if requiredRole != "" && userRole != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			c.Abort()
			return
		}

		// Tiếp tục xử lý request
		c.Next()
	}
}

// Forward request từ API Gateway đến microservice
func forwardRequest(c *gin.Context, targetURL string) {
	var reqBody io.Reader = nil
	if c.Request.Body != nil && (c.Request.Method == "POST" || c.Request.Method == "PUT") {
		reqBody = c.Request.Body
	}

	req, err := http.NewRequest(c.Request.Method, targetURL, reqBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}
	req.Header = c.Request.Header

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward request"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
}

func main() {
	// r := gin.Default()

	// Route đăng nhập để lấy JWT token
	// r.POST("/login", loginHandler)

	// // Áp dụng middleware authentication cho các route còn lại
	// authGroup := r.Group("/")
	// authGroup.Use(authMiddleware())

	// // Route đến User Service
	// authGroup.Any("/users/*any", func(c *gin.Context) {
	// 	targetURL := "http://localhost:5001" + c.Request.URL.Path
	// 	forwardRequest(c, targetURL)
	// })

	// authGroup.Any("/users", func(c *gin.Context) {
	// 	targetURL := "http://localhost:5001/users"
	// 	forwardRequest(c, targetURL)
	// })

	// // Route đến Order Service
	// authGroup.Any("/orders/*any", func(c *gin.Context) {
	// 	targetURL := "http://localhost:5002" + c.Request.URL.Path
	// 	forwardRequest(c, targetURL)
	// })

	// authGroup.Any("/orders", func(c *gin.Context) {
	// 	targetURL := "http://localhost:5002/orders"
	// 	forwardRequest(c, targetURL)
	// })

	// // Route đến Product Service
	// authGroup.Any("/products/*any", func(c *gin.Context) {
	// 	targetURL := "http://localhost:5003" + c.Request.URL.Path
	// 	forwardRequest(c, targetURL)
	// })

	// authGroup.Any("/products", func(c *gin.Context) {
	// 	targetURL := "http://localhost:5003/products"
	// 	forwardRequest(c, targetURL)
	// })

	// // User Service (chỉ admin mới được POST, PUT, DELETE)
	// r.GET("/users/*any", authMiddleware(""), func(c *gin.Context) {
	// 	targetURL := "http://localhost:5001" + c.Request.URL.Path
	// 	forwardRequest(c, targetURL)
	// })

	// r.GET("/users", authMiddleware(""), func(c *gin.Context) {
	// 	targetURL := "http://localhost:5001/users"
	// 	forwardRequest(c, targetURL)
	// })

	// r.POST("/users", authMiddleware("admin"), func(c *gin.Context) {
	// 	targetURL := "http://localhost:5001" + c.Request.URL.Path
	// 	forwardRequest(c, targetURL)
	// })

	// r.PUT("/users/:id", authMiddleware("admin"), func(c *gin.Context) {
	// 	targetURL := "http://localhost:5001" + c.Request.URL.Path
	// 	forwardRequest(c, targetURL)
	// })

	// r.DELETE("/users/:id", authMiddleware("admin"), func(c *gin.Context) {
	// 	targetURL := "http://localhost:5001" + c.Request.URL.Path
	// 	forwardRequest(c, targetURL)
	// })

	r := gin.Default()

	// API Login
	r.POST("/login", loginHandler)

	// Route đến User Service (có auth)
	r.Any("/users/*any", authMiddleware("user"), func(c *gin.Context) {
		targetURL := "http://localhost:5001" + c.Request.URL.Path
		forwardRequest(c, targetURL)
	})

	r.Any("/users", authMiddleware("user"), func(c *gin.Context) {
		targetURL := "http://localhost:5001/users"
		forwardRequest(c, targetURL)
	})

	// Route đến Order Service (chỉ admin mới được truy cập)
	r.Any("/orders/*any", authMiddleware("admin"), func(c *gin.Context) {
		targetURL := "http://localhost:5002" + c.Request.URL.Path
		forwardRequest(c, targetURL)
	})

	r.Any("/orders", authMiddleware("admin"), func(c *gin.Context) {
		targetURL := "http://localhost:5002/orders"
		forwardRequest(c, targetURL)
	})

	fmt.Println("API Gateway running on port 8080...")
	r.Run(":8080")
}
