package routes

import (
	"github.com/febriandani/backend-financial-api/domain/model/general"
	"github.com/febriandani/backend-financial-api/handler/api"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func GetCoreEndpoint(conf *general.AppService, handler api.Handler, log *logrus.Logger) *mux.Router {
	parentRoute := mux.NewRouter()

	jwtRoute := parentRoute.PathPrefix("").Subrouter()
	freeRoute := parentRoute.PathPrefix("").Subrouter()

	// Middleware
	if conf.Authorization.JWT.IsActive {
		log.Info("JWT token is active")
		jwtRoute.Use(handler.Token.JWTValidator)
	}

	// Get Endpoint.
	getV1(freeRoute, jwtRoute, conf, handler)

	return parentRoute
}
