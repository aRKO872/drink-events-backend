package routers

import (
	auth_controller "github.com/drink-events-backend/cmd/controllers/auth-controllers"
	geo_controllers "github.com/drink-events-backend/cmd/controllers/geo-controllers"
	user_controllers "github.com/drink-events-backend/cmd/controllers/user-controllers"
	middlewares "github.com/drink-events-backend/cmd/middleware"
	"github.com/drink-events-backend/cmd/worker-pool"
	"github.com/drink-events-backend/pkg/websockets"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	// Middlewares
	r.Use(middlewares.LoadDatabase)
	r.Use(middlewares.EnableCors)

	go worker.WorkerPool()

	r.Use(middlewares.AllowRoutesMiddleware(middlewares.VerifyToken(), "/geo", "/user", "/ws"))

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

	// GEOLocation Endpoint
	r.POST("/geo/get-nearby-users", geo_controllers.GetNearbyUsers)

	// User Profile Endpoint
	r.POST("/user/change-profile-picture", user_controllers.ChangeProfilePicture)
	r.POST("/user/pair-request-sent-to-me", user_controllers.PairRequestsSentToMe)
	r.POST("/user/pair-request-sent-by-me", user_controllers.PairRequestsSentByMe)
	r.POST("/user/set-search-radius", user_controllers.SetSearchRadius)
	r.POST("/user/delete-pair-request", user_controllers.DeletePairRequest)

	m := websockets.NewWSManager()
	//Websockets
	r.GET("/ws", m.ServeWS)

	// User Websockets Merger
	r.POST("/user/accept-pair-request", m.AcceptPairRequest)

	return r
}
