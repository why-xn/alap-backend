package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/why-xn/alap-backend/pkg/api"
	"github.com/why-xn/alap-backend/pkg/config"
	"github.com/why-xn/alap-backend/pkg/core/context"
	"github.com/why-xn/alap-backend/pkg/db"
	"github.com/why-xn/alap-backend/pkg/db/collection"
	"github.com/why-xn/alap-backend/pkg/db/model"
	"github.com/why-xn/alap-backend/pkg/enum"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"reflect"
)

type UserServiceInterface interface {
	Get(ctx *gin.Context, params *GetUserInputParams) (interface{}, error, int)
	GetList(ctx *gin.Context, params *GetUserListInputParams) (interface{}, error, int)
}

type userService struct{}

var us userService

func UserService() *userService {
	return &us
}

// ---- Get User ---- //

type GetUserInputParams struct {
	UserId string

	output interface{}
	BaseInternalParams
}

func (svc *userService) Get(ctx *gin.Context, p *GetUserInputParams) (interface{}, error, int) {
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

func (p *GetUserInputParams) CheckPermission(ctx *gin.Context) error {
	p.requester = context.GetRequesterFromContext(ctx)
	if p.requester == nil {
		return errors.New("unauthorized request")
	}

	if !p.requester.ResourceAccess[config.KeycloakClientId].HasRole(enum.RoleUser) {
		return errors.New("permission denied")
	}

	return nil
}

func (p *GetUserInputParams) Validate(ctx *gin.Context) error {

	return nil
}

func FetchUser(userId string) (*model.User, error) {
	uid, _ := primitive.ObjectIDFromHex(userId)
	filter := bson.M{"_id": uid, "status": enum.StatusValid}
	res := db.GetDbManager().FindOne(collection.User, filter, reflect.TypeOf(model.User{}))
	var user model.User
	if res != nil {
		user = *res.(*model.User)
		return &user, nil
	} else {
		return nil, errors.New("user not found")
	}
}

func FetchUserByEmail(email string) (*model.User, error) {
	filter := bson.M{"email": email, "status": enum.StatusValid}
	res := db.GetDbManager().FindOne(collection.User, filter, reflect.TypeOf(model.User{}))
	var user model.User
	if res != nil {
		user = *res.(*model.User)
		return &user, nil
	} else {
		return nil, errors.New("user not found")
	}
}

func (p *GetUserInputParams) Process(ctx *gin.Context) (error, int) {
	user, err := FetchUser(p.UserId)
	if err != nil {
		return err, http.StatusNotFound
	}
	p.output = user

	return nil, http.StatusOK
}

func (p *GetUserInputParams) PostProcess(ctx *gin.Context) error {
	return nil
}

// ---- Get All User ---- //

type GetUserListInputParams struct {
	Search *string

	output interface{}
	BaseInternalParams
}

func (svc *userService) GetList(ctx *gin.Context, p *GetUserListInputParams) (interface{}, error, int) {
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

func (p *GetUserListInputParams) CheckPermission(ctx *gin.Context) error {
	p.requester = context.GetRequesterFromContext(ctx)
	if p.requester == nil {
		return errors.New("unauthorized request")
	}

	if !p.requester.ResourceAccess[config.KeycloakClientId].HasRole(enum.RoleUser) {
		return errors.New("permission denied")
	}

	return nil
}

func (p *GetUserListInputParams) Validate(ctx *gin.Context) error {
	return nil
}

func (p *GetUserListInputParams) Process(ctx *gin.Context) (error, int) {
	filter := bson.M{
		"status": enum.StatusValid,
	}

	var userList []model.User
	res := db.GetDbManager().FindAll(collection.User, reflect.TypeOf(model.User{}), filter, nil, -1, 100)
	if res != nil {
		userList = db.ConvertToUserArray(res)
	}

	p.output = userList

	return nil, http.StatusOK
}

func (p *GetUserListInputParams) PostProcess(ctx *gin.Context) error {
	return nil
}
