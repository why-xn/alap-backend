package v1

import (
	"github.com/why-xn/alap-backend/pkg/core/keycloak"
)

type BaseInternalParams struct {
	requester *keycloak.SsoUserDTO
}
