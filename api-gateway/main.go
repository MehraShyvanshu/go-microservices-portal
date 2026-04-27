package main

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func proxyRequest(c *gin.Context, method string, url string) {
	bodyBytes, _ := io.ReadAll(c.Request.Body)

	req, _ := http.NewRequest(method, url, bytes.NewBuffer(bodyBytes))
	req.Header = c.Request.Header

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		c.JSON(500, gin.H{"error": "service unavailable"})
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", respBody)
}

func main() {
	r := gin.Default()

	// PUBLIC
	r.POST("/login", func(c *gin.Context) {
		proxyRequest(c, "POST", "http://localhost:8001/login")
	})

	r.POST("/register", func(c *gin.Context) {
		proxyRequest(c, "POST", "http://localhost:8001/register")
	})

	// PROTECTED
	api := r.Group("/api")
	api.Use(AuthMiddleware())

	// 🔒 ADMIN ONLY
	api.GET("/users",
		RoleMiddleware("admin"),
		func(c *gin.Context) {
			proxyRequest(c, "GET", "http://localhost:8004/users")
		})

	// 🔓 USER + HR + ADMIN
	api.GET("/tasks",
		RoleMiddleware("admin", "user", "hr"),
		func(c *gin.Context) {
			proxyRequest(c, "GET", "http://localhost:8003/tasks")
		})

	api.POST("/tasks",
		RoleMiddleware("admin", "user", "hr"),
		func(c *gin.Context) {
			proxyRequest(c, "POST", "http://localhost:8003/tasks")
		})

	api.PUT("/role",
		RoleMiddleware("admin"),
		func(c *gin.Context) {
			proxyRequest(c, "PUT", "http://localhost:8001/role")
		})

	r.Run(":8080")
}
