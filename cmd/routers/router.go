package routers

import (
	auth_controller "github.com/drink-events-backend/cmd/controllers/auth-controllers"
	middlewares "github.com/drink-events-backend/cmd/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	// Middlewares
	r.Use(middlewares.LoadEnv)
	r.Use(middlewares.LoadDatabase)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Auth Endpoints
	r.POST("/auth/verify-email", auth_controller.VerifyEmail)
	r.POST("/auth/verify-otp", auth_controller.VerifyOTP)
	r.POST("/auth/resend-otp-verify", auth_controller.ResendEmailOTP)
	r.POST("/auth/signup", auth_controller.SignUp)
	r.POST("/auth/login", auth_controller.LogIn)

	return r
}
