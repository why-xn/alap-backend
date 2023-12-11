package keycloak

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Nerzal/gocloak/v13"
	"github.com/golang-jwt/jwt/v4"
	config "github.com/why-xn/alap-backend/pkg/config"
	"github.com/why-xn/alap-backend/pkg/core/log"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var kClient *gocloak.GoCloak
var adminKClient *gocloak.GoCloak
var token *gocloak.JWT

type AccessDTO struct {
	Roles []string `json:"roles"`
}

func (access AccessDTO) HasRole(role string) bool {
	if len(access.Roles) == 0 {
		return false
	}

	for _, r := range access.Roles {
		if r == role {
			return true
		}
	}
	return false
}

type SsoUserDTO struct {
	Azp            *string              `json:"azp,omitempty"`
	Name           string               `json:"name"`
	Email          string               `json:"email"`
	EmailVerified  bool                 `json:"email_verified"`
	Username       string               `json:"preferred_username"`
	RealmAccess    *AccessDTO           `json:"realm_access,omitempty"`
	ResourceAccess map[string]AccessDTO `json:"resource_access,omitempty"`
	Roles          []string             `json:"roles"`
	SessionState   string               `json:"session_state"`
	Exp            float64              `json:"exp"`
}

func InitKeycloakClient() {
	kClient = gocloak.NewClient(config.KeycloakHost)
}

func AuthenticateUser(authorizationCode string, state string, redirectUri string) (error, *gocloak.JWT) {
	tokenApiUrl := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", config.KeycloakHost, config.KeycloakRealm)

	formData := url.Values{}
	formData.Set("client_id", config.KeycloakClientId)
	formData.Set("client_secret", config.KeycloakClientSecret)
	formData.Set("scope", "openid email profile")
	formData.Set("redirect_uri", redirectUri)
	formData.Set("state", state)
	formData.Set("grant_type", "authorization_code")
	formData.Set("code", authorizationCode)

	encodedData := formData.Encode()

	req, err := http.NewRequest("POST", tokenApiUrl, strings.NewReader(encodedData))
	if err != nil {
		return err, nil
	}

	// Set the Content-Type header to application/x-www-form-urlencoded
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return err, nil
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err, nil
	}

	var authResponse gocloak.JWT
	err = json.Unmarshal(body, &authResponse)
	if err != nil {
		log.Logger.Fatalw("failed to unmarshal auth response", err, err.Error())
		return err, nil
	} else {
		log.Logger.Debugw("Auth response", "accessToken", authResponse.AccessToken, "expiresIn", authResponse.ExpiresIn)
	}

	return nil, &authResponse
}

func LogoutUser(ctx context.Context, refreshToken string) error {
	return kClient.Logout(ctx, config.KeycloakClientId, config.KeycloakClientSecret, config.KeycloakRealm, refreshToken)
}

func RefreshToken(ctx context.Context, refreshToken string) (*gocloak.JWT, error) {
	return kClient.RefreshToken(ctx, refreshToken, config.KeycloakClientId, config.KeycloakClientSecret, config.KeycloakRealm)
}

func DecodeAccessToken(accessToken string) (error, jwt.MapClaims) {
	_, mapClaims, err := kClient.DecodeAccessToken(context.Background(), accessToken, config.KeycloakRealm)
	if err != nil {
		log.Logger.Debugw("failed to decode access token", err, err.Error())
		return err, nil
	} else {
		//log.Logger.Infow("decoded access token", "claims", mapClaims)
		return nil, *mapClaims
	}
}

func MapClaimsToSsoUserDTO(claims jwt.MapClaims) *SsoUserDTO {
	if claims == nil {
		return nil
	}

	claimsData, err := json.Marshal(claims)
	if err != nil {
		return nil
	}

	var ssoUserDto SsoUserDTO
	err = json.Unmarshal(claimsData, &ssoUserDto)
	if err != nil {
		return nil
	}
	return &ssoUserDto
}

func HasAccessToThisApp(claimsDTO *SsoUserDTO) bool {
	if claimsDTO.ResourceAccess != nil {
		if _, ok := claimsDTO.ResourceAccess[config.KeycloakClientId]; ok {
			return true
		}
	}
	return false
}

func IsTokenExpired(accessToken string) bool {
	_, claims := DecodeAccessToken(accessToken)
	if claims != nil {
		if exp, ok := claims["exp"].(float64); ok {
			expirationTime := time.Unix(int64(exp), 0)
			return time.Now().After(expirationTime)
		}
	}
	return true
}
