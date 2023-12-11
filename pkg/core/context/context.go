package context

import (
	"github.com/gin-gonic/gin"
	"github.com/why-xn/alap-backend/pkg/core/keycloak"
)

func AddAccessTokenToContext(c *gin.Context, accessToken interface{}) {
	c.Set("AccessToken", accessToken)
}

func AddRequesterToContext(c *gin.Context, user *keycloak.SsoUserDTO) {
	c.Set("User", user)
}

func AddWsIdToContext(c *gin.Context, wsId string) {
	c.Set("WsId", wsId)
}

func GetWsIdFromContext(c *gin.Context) string {
	if val, ok := c.Get("WsId"); ok {
		return val.(string)
	}
	return ""
}

func GetRequesterFromContext(c *gin.Context) *keycloak.SsoUserDTO {
	if val, ok := c.Get("User"); ok {
		if val == nil {
			return nil
		}
		user, ok := val.(*keycloak.SsoUserDTO)
		if !ok {
			return nil
		}
		return user
	}
	return nil
}
