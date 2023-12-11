package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/why-xn/alap-backend/pkg/api"
	"github.com/why-xn/alap-backend/pkg/core/context"
	"github.com/why-xn/alap-backend/pkg/core/keycloak"
	"github.com/why-xn/alap-backend/pkg/core/log"
	"github.com/why-xn/alap-backend/pkg/db"
	"github.com/why-xn/alap-backend/pkg/db/collection"
	"github.com/why-xn/alap-backend/pkg/db/model"
	"github.com/why-xn/alap-backend/pkg/enum"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"reflect"
	"time"
)

type AuthServiceInterface interface {
	Login(ctx *gin.Context, params *LoginInputParams) (interface{}, error, int)
	Logout(ctx *gin.Context, params *LogoutInputParams) (interface{}, error, int)
	RefreshToken(ctx *gin.Context, params *RefreshTokenInputParams) (interface{}, error, int)
}

type authService struct{}

var as authService

func AuthService() *authService {
	return &as
}

// ---- Login ---- //

type LoginInputParams struct {
	AuthorizationCode string `json:"authorizationCode"`
	State             string `json:"state"`
	RedirectUri       string `json:"redirectUri"`

	output interface{}
	BaseInternalParams
}

type LoginOutputDTO struct {
	AccessToken      string `json:"accessToken"`
	RefreshToken     string `json:"refreshToken"`
	SessionState     string `json:"sessionState"`
	ExpiresIn        int    `json:"expireIn"`
	RefreshExpiresIn int    `json:"refreshExpiresIn"`
	TokenType        string `json:"tokenType"`

	UserInfo model.User `json:"userInfo"`
}

func (svc *authService) Login(ctx *gin.Context, p *LoginInputParams) (interface{}, error, int) {
	err := p.CheckPermission(ctx)
	if err != nil {
		return api.ErrorResponse(err, http.StatusUnauthorized)
	}

	err = p.Validate(ctx)
	if err != nil {
		return api.ErrorResponse(err, http.StatusBadRequest)
	}

	err, code := p.Process(ctx)
	if err != nil {
		return api.ErrorResponse(err, code)
	}

	_ = p.PostProcess(ctx)

	return p.output, nil, code
}

func (p *LoginInputParams) CheckPermission(ctx *gin.Context) error {
	return nil
}

func (p *LoginInputParams) Validate(ctx *gin.Context) error {

	return nil
}

func (p *LoginInputParams) Process(ctx *gin.Context) (error, int) {

	err, token := keycloak.AuthenticateUser(p.AuthorizationCode, p.State, p.RedirectUri)
	if err != nil {
		log.Logger.Errorw("failed to authenticate user", "err", err.Error())
		return errors.New("failed to authenticate user"), http.StatusUnauthorized
	}

	err, claims := keycloak.DecodeAccessToken(token.AccessToken)
	if err != nil || claims == nil {
		return errors.New("invalid authorization code"), http.StatusUnauthorized
	}

	ssoUserDTO := keycloak.MapClaimsToSsoUserDTO(claims)

	// user email has to be verified to login
	if !ssoUserDTO.EmailVerified {
		return errors.New("user not verified"), http.StatusUnauthorized
	}

	// user need to have at-least one role from keycloak client to login
	hasAccessToApp := keycloak.HasAccessToThisApp(ssoUserDTO)
	if !hasAccessToApp {
		return errors.New("user do not have permission to access to this application"), http.StatusUnauthorized
	}

	filter := bson.M{"email": ssoUserDTO.Email, "status": enum.StatusValid}
	res := db.GetDbManager().FindOne(collection.User, filter, reflect.TypeOf(model.User{}))
	var user model.User
	if res == nil {
		user = model.User{
			ID:            primitive.NewObjectID(),
			FullName:      ssoUserDTO.Name,
			Email:         ssoUserDTO.Email,
			Username:      ssoUserDTO.Email,
			Status:        enum.StatusValid,
			CreatedAt:     time.Now(),
			LastUpdatedAt: time.Now(),
		}
		_, err := db.GetDbManager().InsertSingleDocument(collection.User, user)
		if err != nil {
			log.Logger.Errorw("Failed to save user in db", "err", err.Error())
			return errors.New("internal server error"), http.StatusInternalServerError
		}
	} else {
		user = *res.(*model.User)
	}

	out := LoginOutputDTO{
		TokenType:        token.TokenType,
		AccessToken:      token.AccessToken,
		ExpiresIn:        token.ExpiresIn,
		RefreshToken:     token.RefreshToken,
		RefreshExpiresIn: token.RefreshExpiresIn,
		SessionState:     token.SessionState,
		UserInfo:         user,
	}

	p.output = out

	return nil, http.StatusOK
}

func (p *LoginInputParams) PostProcess(ctx *gin.Context) error {
	return nil
}

// ---- Logout ---- //

type LogoutInputParams struct {
	RefreshToken string `json:"refreshToken"`

	output interface{}
	BaseInternalParams
}

func (svc *authService) Logout(ctx *gin.Context, p *LogoutInputParams) (interface{}, error, int) {
	err := p.CheckPermission(ctx)
	if err != nil {
		return api.ErrorResponse(err, http.StatusUnauthorized)
	}

	err = p.Validate(ctx)
	if err != nil {
		return api.ErrorResponse(err, 400)
	}

	err, code := p.Process(ctx)
	if err != nil {
		return api.ErrorResponse(err, code)
	}

	_ = p.PostProcess(ctx)

	return p.output, nil, code
}

func (p *LogoutInputParams) CheckPermission(ctx *gin.Context) error {
	p.requester = context.GetRequesterFromContext(ctx)
	if p.requester == nil {
		return errors.New("unauthorized request")
	}
	return nil
}

func (p *LogoutInputParams) Validate(ctx *gin.Context) error {
	if len(p.RefreshToken) == 0 {
		return errors.New("refresh token is required")
	}
	return nil
}

func (p *LogoutInputParams) Process(ctx *gin.Context) (error, int) {
	err := keycloak.LogoutUser(ctx, p.RefreshToken)
	if err != nil {
		log.Logger.Errorw("Failed to logout user", "err", err.Error())
		return err, http.StatusBadRequest
	}
	return nil, http.StatusOK
}

func (p *LogoutInputParams) PostProcess(ctx *gin.Context) error {
	return nil
}

// ---- Refresh Token ---- //

type RefreshTokenInputParams struct {
	RefreshToken string `json:"refreshToken"`

	output interface{}
	BaseInternalParams
}

func (svc *authService) RefreshToken(ctx *gin.Context, p *RefreshTokenInputParams) (interface{}, error, int) {
	err := p.CheckPermission(ctx)
	if err != nil {
		return api.ErrorResponse(err, http.StatusUnauthorized)
	}

	err = p.Validate(ctx)
	if err != nil {
		return api.ErrorResponse(err, 400)
	}

	err, code := p.Process(ctx)
	if err != nil {
		return api.ErrorResponse(err, code)
	}

	_ = p.PostProcess(ctx)

	return p.output, nil, code
}

func (p *RefreshTokenInputParams) CheckPermission(ctx *gin.Context) error {
	return nil
}

func (p *RefreshTokenInputParams) Validate(ctx *gin.Context) error {
	if len(p.RefreshToken) == 0 {
		return errors.New("refresh token is required")
	}
	return nil
}

func (p *RefreshTokenInputParams) Process(ctx *gin.Context) (error, int) {
	newToken, err := keycloak.RefreshToken(ctx, p.RefreshToken)
	if err != nil {
		log.Logger.Errorw("Failed to refresh token", "err", err.Error())
		return err, http.StatusBadRequest
	}
	p.output = newToken
	return nil, http.StatusOK
}

func (p *RefreshTokenInputParams) PostProcess(ctx *gin.Context) error {
	return nil
}
