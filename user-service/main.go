package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID    int
	Name  string
	Email string
}

var users = []User{
	{1, "John", "john@test.com"},
	{2, "Alice", "alice@test.com"},
}

func main() {
	r := gin.Default()

	r.GET("/users", func(c *gin.Context) {
		c.JSON(http.StatusOK, users)
	})

	r.Run(":8004")
}
