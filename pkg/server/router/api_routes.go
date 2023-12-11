package router

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/why-xn/alap-backend/pkg/api/controller/v1"
)

// @title           Swagger API
// @version         1.0

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey  ApiKeyAuth
// @in header
// @name Authorization

func AddApiRoutes(httpRg *gin.RouterGroup) {
	//httpRg.Use(TokenAuthMiddleware())
	httpRg.POST("api/v1/auth/login", v1.AuthController().Login)
	httpRg.POST("api/v1/auth/logout", v1.AuthController().Logout)
	httpRg.POST("api/v1/auth/refresh-token", v1.AuthController().RefreshToken)
	httpRg.GET("api/v1/auth/who-am-i", v1.AuthController().WhoAmI)

	httpRg.GET("api/v1/users", v1.UserController().GetList)
	httpRg.GET("api/v1/users/:id", v1.UserController().Get)

	httpRg.GET("api/v1/chat-window", v1.ChatWindowController().Get)
}
