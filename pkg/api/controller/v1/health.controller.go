package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HealthControllerInterface interface {
	Check(c *gin.Context)
}

type healthController struct{}

var hc healthController

func HealthController() *healthController {
	return &hc
}

// Check godoc
// @Summary      Check Health
// @Description  Check Health
// @Tags         health
// @Accept       json
// @Produce      json
// @Success      200  {object}  api.ResponseDTO
// @Failure      400  {object}  api.ResponseDTO
// @Router       /health [get]
func (ctrl *healthController) Check(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "i am alive",
	})
}
