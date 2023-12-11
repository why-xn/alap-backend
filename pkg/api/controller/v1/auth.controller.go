package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/why-xn/alap-backend/pkg/api"
	v1 "github.com/why-xn/alap-backend/pkg/api/service/v1"
	"github.com/why-xn/alap-backend/pkg/core/context"
	"github.com/why-xn/alap-backend/pkg/core/log"
)

type AuthControllerInterface interface {
	Login(c *gin.Context)
	Logout(c *gin.Context)
	RefreshToken(c *gin.Context)
}

type authController struct{}

var ac authController

func AuthController() *authController {
	return &ac
}

// Login godoc
// @Summary      Login
// @Description  Login
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  api.ResponseDTO
// @Failure      400  {object}  api.ResponseDTO
// @Router       /v1/auth/login [post]
func (ctrl *authController) Login(ctx *gin.Context) {
	var input v1.LoginInputParams

	err := ctx.BindJSON(&input)
	if err != nil {
		log.Logger.Errorw("Failed to bind JSON", "err", err)
		return
	}

	res, err, code := v1.AuthService().Login(ctx, &input)
	if err != nil {
		log.Logger.Debugw("API Execution failed", "err", err)
		api.SendErrorResponse(ctx, err.Error(), code)
		return
	}

	api.SendResponse(ctx, res)
}

// Logout godoc
// @Summary      Logout
// @Description  Logout
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  api.ResponseDTO
// @Failure      400  {object}  api.ResponseDTO
// @Router       /v1/auth/logout [post]
func (ctrl *authController) Logout(ctx *gin.Context) {
	var input v1.LogoutInputParams

	err := ctx.BindJSON(&input)
	if err != nil {
		log.Logger.Errorw("Failed to bind JSON", "err", err)
		return
	}

	res, err, code := v1.AuthService().Logout(ctx, &input)
	if err != nil {
		log.Logger.Debugw("API Execution failed", "err", err)
		api.SendErrorResponse(ctx, err.Error(), code)
		return
	}

	api.SendResponse(ctx, res)
}

// RefreshToken godoc
// @Summary      Refresh Token
// @Description  Refresh Token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  api.ResponseDTO
// @Failure      400  {object}  api.ResponseDTO
// @Router       /v1/auth/refresh-token [post]
func (ctrl *authController) RefreshToken(ctx *gin.Context) {
	var input v1.RefreshTokenInputParams

	err := ctx.BindJSON(&input)
	if err != nil {
		log.Logger.Errorw("Failed to bind JSON", "err", err)
		return
	}

	res, err, code := v1.AuthService().RefreshToken(ctx, &input)
	if err != nil {
		log.Logger.Debugw("API Execution failed", "err", err)
		api.SendErrorResponse(ctx, err.Error(), code)
		return
	}

	api.SendResponse(ctx, res)
}

type WhoAmIInputParams struct {
	AccessToken string `json:"accessToken"`
}

// WhoAmI godoc
// @Summary      WhoAmI
// @Description  WhoAmI
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  api.ResponseDTO
// @Failure      400  {object}  api.ResponseDTO
// @Router       /v1/auth/who-am-i [get]
func (ctrl *authController) WhoAmI(ctx *gin.Context) {
	user := context.GetRequesterFromContext(ctx)
	if user == nil {
		api.SendErrorResponse(ctx, "unauthorized request", 401)
		return
	}
	api.SendResponse(ctx, user)
}
