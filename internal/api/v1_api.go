package api

import (
	"net/http"

	"github.com/config-source/cdb/internal/apiutils"
	"github.com/config-source/cdb/internal/middleware"
	"github.com/config-source/cdb/pkg/auth"
	"github.com/config-source/cdb/pkg/configkeys"
	"github.com/config-source/cdb/pkg/configvalues"
	"github.com/config-source/cdb/pkg/environments"
	"github.com/config-source/cdb/pkg/services"
	"github.com/rs/zerolog"
)

type V1 struct {
	log             zerolog.Logger
	tokenSigningKey []byte

	userService        *auth.UserService
	configValueService *configvalues.Service
	envService         *environments.Service
	svcService         *services.ServiceService
	configKeyService   *configkeys.Service
}

func NewV1(
	log zerolog.Logger,
	tokenSigningKey []byte,
	userService *auth.UserService,
	configValueService *configvalues.Service,
	envService *environments.Service,
	configKeyService *configkeys.Service,
	svcService *services.ServiceService,
) (*V1, http.Handler) {
	api := &V1{
		log:             log,
		tokenSigningKey: tokenSigningKey,

		configValueService: configValueService,
		envService:         envService,
		configKeyService:   configKeyService,
		userService:        userService,
		svcService:         svcService,
	}

	// v1 routes
	v1Mux := http.NewServeMux()

	v1Mux.HandleFunc("GET /api/v1/environments/{serviceName}/by-name/{name}", api.GetEnvironmentByName)
	v1Mux.HandleFunc("GET /api/v1/environments/by-id/{id}", api.GetEnvironmentByID)
	v1Mux.HandleFunc("GET /api/v1/environments/tree", api.GetEnvironmentTree)
	v1Mux.HandleFunc("GET /api/v1/environments", api.ListEnvironments)
	v1Mux.HandleFunc("POST /api/v1/environments", api.CreateEnvironment)
	v1Mux.HandleFunc("PUT /api/v1/environments/{id}", api.UpdateEnvironment)
	v1Mux.HandleFunc("DELETE /api/v1/environments/{id}", api.DeleteEnvironment)

	v1Mux.HandleFunc("GET /api/v1/services/by-name/{name}", api.GetServiceByName)
	v1Mux.HandleFunc("GET /api/v1/services/by-id/{id}", api.GetServiceByID)
	v1Mux.HandleFunc("GET /api/v1/services", api.ListServices)
	v1Mux.HandleFunc("POST /api/v1/services", api.CreateService)

	v1Mux.HandleFunc("POST /api/v1/config-keys", api.CreateConfigKey)
	v1Mux.HandleFunc("GET /api/v1/config-keys", api.ListConfigKeys)
	v1Mux.HandleFunc("GET /api/v1/config-keys/by-id/{id}", api.GetConfigKeyByID)
	v1Mux.HandleFunc("GET /api/v1/config-keys/{serviceName}/by-name/{name}", api.GetConfigKeyByName)

	v1Mux.HandleFunc("POST /api/v1/config-values", api.CreateConfigValue)
	v1Mux.HandleFunc("GET /api/v1/config-values/{environment}/{key}", api.GetConfigurationValue)
	v1Mux.HandleFunc("POST /api/v1/config-values/{environment}/{key}", api.SetConfigurationValue)
	v1Mux.HandleFunc("GET /api/v1/config-values/{environment}", api.GetConfiguration)

	v1Mux.HandleFunc("GET /api/v1/users/me", api.GetLoggedInUser)
	v1Mux.HandleFunc("POST /api/v1/auth/api-tokens", api.IssueAPIToken)
	v1Mux.HandleFunc("GET /api/v1/auth/api-tokens", api.ListAPITokens)

	apiMux := http.NewServeMux()
	apiMux.HandleFunc("DELETE /api/v1/auth/logout", api.Logout)
	apiMux.HandleFunc("POST /api/v1/auth/login", api.Login)
	apiMux.HandleFunc("POST /api/v1/auth/register", api.Register)
	apiMux.Handle("/api/v1/", middleware.AuthenticationRequired(log, userService, tokenSigningKey, v1Mux))

	return api, apiMux
}

func (a *V1) sendJson(w http.ResponseWriter, payload interface{}) {
	apiutils.SendJSON(a.log, w, payload)
}

func (a *V1) sendErr(w http.ResponseWriter, err error) {
	apiutils.SendErr(a.log, w, err)
}
