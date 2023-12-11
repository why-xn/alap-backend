package router

import (
	"github.com/gin-gonic/gin"
	"github.com/why-xn/alap-backend/pkg/api/controller/v1"
)

// @title           Swagger API
// @version         1.0

// @host      localhost:8080
// @BasePath  /

func AddGeneralRoutes(httpRg *gin.RouterGroup) {
	httpRg.GET("health", v1.HealthController().Check)
}
