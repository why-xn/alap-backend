package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/why-xn/alap-backend/pkg/api"
	v1 "github.com/why-xn/alap-backend/pkg/api/service/v1"
	"github.com/why-xn/alap-backend/pkg/core/log"
)

type UserControllerInterface interface {
	Get(c *gin.Context)
	GetList(c *gin.Context)
}

type userController struct{}

var uc userController

func UserController() *userController {
	return &uc
}

// Get godoc
// @Summary      Create User
// @Description  Create User
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200  {object}  api.ResponseDTO
// @Failure      400  {object}  api.ResponseDTO
// @Router       /v1/user/:id [get]
func (ctrl *userController) Get(ctx *gin.Context) {
	input := new(v1.GetUserInputParams)
	input.UserId = ctx.Param("id")

	res, err, code := v1.UserService().Get(ctx, input)
	if err != nil {
		log.Logger.Debugw("API Execution failed", "err", err)
		api.SendErrorResponse(ctx, err.Error(), code)
		return
	}

	api.SendResponse(ctx, res)
}

// GetList godoc
// @Summary      Get User List
// @Description  Get User List
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200  {object}  api.ResponseDTO
// @Failure      400  {object}  api.ResponseDTO
// @Router       /v1/user [get]
func (ctrl *userController) GetList(ctx *gin.Context) {
	input := new(v1.GetUserListInputParams)
	search := ctx.Query("search")

	if len(search) > 0 {
		input.Search = &search
	}

	res, err, code := v1.UserService().GetList(ctx, input)
	if err != nil {
		log.Logger.Debugw("API Execution failed", "err", err)
		api.SendErrorResponse(ctx, err.Error(), code)
		return
	}

	api.SendResponse(ctx, res)
}
