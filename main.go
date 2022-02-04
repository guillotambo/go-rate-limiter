package main

import (
	"github.com/gin-gonic/gin"
	"lemon_cash/controller"
	"lemon_cash/middleware"
)

func main() {
	router := gin.New()
	rateLimiter := middleware.NewRateLimiter(5, 10)

	router.Group("/message", rateLimiter.Accept).
		GET("/", controller.GetFOAASMessage)

	router.Run(":8080")
}
