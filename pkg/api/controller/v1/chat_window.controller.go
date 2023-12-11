package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/why-xn/alap-backend/pkg/api"
	v1 "github.com/why-xn/alap-backend/pkg/api/service/v1"
	"github.com/why-xn/alap-backend/pkg/core/log"
)

type ChatWindowControllerInterface interface {
	Get(c *gin.Context)
}

type chatWindowController struct{}

var cwc chatWindowController

func ChatWindowController() *chatWindowController {
	return &cwc
}

// Get godoc
// @Summary      Create User
// @Description  Create User
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200  {object}  api.ResponseDTO
// @Failure      400  {object}  api.ResponseDTO
// @Router       /v1/chat-window/:id [get]
func (ctrl *chatWindowController) Get(ctx *gin.Context) {
	input := new(v1.GetChatWindowInputParams)
	//input.ChatWindowId = ctx.Param("id")
	input.ToUser = ctx.Query("to")

	res, err, code := v1.ChatWindowService().Get(ctx, input)
	if err != nil {
		log.Logger.Debugw("API Execution failed", "err", err)
		api.SendErrorResponse(ctx, err.Error(), code)
		return
	}

	api.SendResponse(ctx, res)
}
