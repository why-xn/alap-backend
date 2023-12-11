package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/why-xn/alap-backend/pkg/api"
	"github.com/why-xn/alap-backend/pkg/core/context"
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

type ChatWindowServiceInterface interface {
	Get(ctx *gin.Context, params *GetChatWindowInputParams) (interface{}, error, int)
}

type chatWindowService struct{}

var cws chatWindowService

func ChatWindowService() *chatWindowService {
	return &cws
}

// ---- Get ChatWindow ---- //

type GetChatWindowInputParams struct {
	ChatWindowId string
	ToUser       string

	output interface{}
	BaseInternalParams
}

func (svc *chatWindowService) Get(ctx *gin.Context, p *GetChatWindowInputParams) (interface{}, error, int) {
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

func (p *GetChatWindowInputParams) CheckPermission(ctx *gin.Context) error {
	p.requester = context.GetRequesterFromContext(ctx)
	if p.requester == nil {
		return errors.New("unauthorized request")
	}
	return nil
}

func (p *GetChatWindowInputParams) Validate(ctx *gin.Context) error {

	return nil
}

func (p *GetChatWindowInputParams) Process(ctx *gin.Context) (error, int) {
	toUser, err := FetchUser(p.ToUser)
	if err != nil {
		return err, http.StatusInternalServerError
	}

	senderUser, err := FetchUserByEmail(p.requester.Email)
	if err != nil {
		return err, http.StatusInternalServerError
	}

	filter := bson.M{"users": p.ToUser, "status": enum.StatusValid}
	res := db.GetDbManager().FindOne(collection.ChatWindow, filter, reflect.TypeOf(model.ChatWindow{}))
	var chatWindow model.ChatWindow
	if res != nil {
		chatWindow = *res.(*model.ChatWindow)
	} else {
		chatWindow = model.ChatWindow{
			ID: primitive.NewObjectID(),
			Users: []string{
				senderUser.ID.Hex(),
				toUser.ID.Hex(),
			},
			Status:        enum.StatusValid,
			CreatedAt:     time.Now(),
			LastUpdatedAt: time.Now(),
		}

		_, err := db.GetDbManager().InsertSingleDocument(collection.ChatWindow, chatWindow)
		if err != nil {
			log.Logger.Errorw("Failed to save user in db", "err", err.Error())
			return errors.New("internal server error"), http.StatusInternalServerError
		}
	}

	p.output = chatWindow

	return nil, http.StatusOK
}

func (p *GetChatWindowInputParams) PostProcess(ctx *gin.Context) error {
	return nil
}
