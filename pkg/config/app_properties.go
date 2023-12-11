package config

import (
	"github.com/joho/godotenv"
	"github.com/why-xn/alap-backend/pkg/core/log"
	"github.com/why-xn/alap-backend/pkg/types"
	"os"
)

var RunMode types.RunMode
var ServerPort string
var KeycloakHost string
var KeycloakRealm string
var KeycloakClientId string
var KeycloakClientSecret string
var DatabaseConnectionString string
var DatabaseName string

func InitEnvironmentVariables() {
	RunMode = types.RunMode(os.Getenv("RUN_MODE"))
	if RunMode == "" {
		RunMode = Local
	}

	log.Logger.Info("Run Mode: ", RunMode)

	if RunMode != K8 {
		err := godotenv.Load()
		if err != nil {
			log.Logger.Errorw("Failed to load environment variables from .env file", "err", err.Error())
			return
		}
	}

	ServerPort = os.Getenv("SERVER_PORT")

	KeycloakHost = os.Getenv("KEYCLOAK_HOST")
	KeycloakRealm = os.Getenv("KEYCLOAK_REALM")
	KeycloakClientId = os.Getenv("KEYCLOAK_CLIENT_ID")
	KeycloakClientSecret = os.Getenv("KEYCLOAK_CLIENT_SECRET")

	DatabaseConnectionString = os.Getenv("DATABASE_CONNECTION_STRING")
	DatabaseName = os.Getenv("DATABASE_NAME")
}
