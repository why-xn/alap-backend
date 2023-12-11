package router

import (
	"github.com/gin-gonic/gin"
	"github.com/why-xn/alap-backend/pkg/core/context"
	"github.com/why-xn/alap-backend/pkg/core/keycloak"
	"github.com/why-xn/alap-backend/pkg/core/log"
	"net/http"
)

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if accessToken := c.GetHeader("Authorization"); len(accessToken) > 0 {
			err, claimMap := keycloak.DecodeAccessToken(accessToken)
			if err != nil {
				log.Logger.Debugw("Error occurred while decoding access token", "err", err.Error())
				c.JSON(http.StatusUnauthorized, logError("invalid token"))
				c.Abort()
			}

			ssoUserDTO := keycloak.MapClaimsToSsoUserDTO(claimMap)
			context.AddAccessTokenToContext(c, accessToken)
			context.AddRequesterToContext(c, ssoUserDTO)
		}
		c.Next()
	}
}

func logError(errMsg string) gin.H {
	return gin.H{"msg": errMsg}
}
