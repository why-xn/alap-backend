package server

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	docs "github.com/why-xn/alap-backend/docs"
	v1 "github.com/why-xn/alap-backend/pkg/api/service/v1"
	"github.com/why-xn/alap-backend/pkg/config"
	"github.com/why-xn/alap-backend/pkg/core/context"
	"github.com/why-xn/alap-backend/pkg/core/keycloak"
	"github.com/why-xn/alap-backend/pkg/core/log"
	"github.com/why-xn/alap-backend/pkg/core/messaging"
	"github.com/why-xn/alap-backend/pkg/dto"
	"github.com/why-xn/alap-backend/pkg/server/router"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	// Solve cross-domain problems
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

func ws(c *gin.Context) {
	//Upgrade get request to webSocket protocol
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Logger.Info("upgrade:", err)
		return
	}
	defer ws.Close()

	token, ok := c.GetQuery("accessToken")
	if ok {
		log.Logger.Debug("Token ", token)
	}

	err, claimMap := keycloak.DecodeAccessToken(token)
	if err != nil {
		log.Logger.Errorw("Error occurred while decoding access token", "err", err.Error())
		return
	}

	ssoUserDTO := keycloak.MapClaimsToSsoUserDTO(claimMap)
	context.AddAccessTokenToContext(c, token)
	context.AddRequesterToContext(c, ssoUserDTO)

	user, err := v1.FetchUserByEmail(ssoUserDTO.Email)
	if err != nil {
		log.Logger.Errorw("Error occurred fetching use", "err", err.Error())
		return
	}

	wsId := uuid.NewString()
	context.AddWsIdToContext(c, wsId)
	messaging.AddToWsMap(wsId, ws)
	messaging.AddWsToUser(user.ID.Hex(), wsId)

	defer messaging.RemoveFromWsMap(wsId)
	defer messaging.RemoveWsFromUser(user.ID.Hex(), wsId)

	log.Logger.Infow("User connected", "user", ssoUserDTO.Username, "ws", wsId)

	var wsInMsg dto.WsInMessageDTO

	for {
		//read data from ws
		mt, message, err := ws.ReadMessage()
		if err != nil {
			log.Logger.Info("read:", err)
			break
		}

		err = json.Unmarshal(message, &wsInMsg)
		if err != nil {
			log.Logger.Errorw("failed to parse incoming msg", "err", err.Error())
		}

		wsInMsg.MessageType = mt

		messaging.ProcessIncomingMsg(wsInMsg)

		//write ws data
		/*err = ws.WriteMessage(mt, message)
		if err != nil {
			log.Logger.Info("write:", err)
			break
		}*/
	}

	log.Logger.Infow("Disconnected", "ws", wsId)

}

func Start() {
	r := gin.Default()

	r.Use(router.TokenAuthMiddleware())

	// Setup CORS Config
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.MaxAge = time.Second * 3600
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Content-Type")
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization")
	corsConfig.ExposeHeaders = append(corsConfig.ExposeHeaders, "Content-Length")
	r.Use(cors.New(corsConfig))

	docs.SwaggerInfo.BasePath = "/"

	// Setting API Base Path for HTTP APIs
	httpRouter := r.Group("/")

	// Setting up all Http Routes
	router.AddGeneralRoutes(httpRouter)
	router.AddApiRoutes(httpRouter)

	// Setup Swagger Route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/ws", ws)

	log.Logger.Infof("Starting Web Server in port %s", config.ServerPort)
	err := r.Run(fmt.Sprintf(":%s", config.ServerPort)) // listen and serve on 0.0.0.0:PORT
	if err != nil {
		log.Logger.Errorw("Failed to start Web Server", "err", err.Error())
	}
}
