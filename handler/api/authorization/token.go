package authorization

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	cg "github.com/febriandani/backend-financial-api/domain/constants/general"
	dg "github.com/febriandani/backend-financial-api/domain/model/general"
	"github.com/febriandani/backend-financial-api/domain/utils"
	"github.com/sirupsen/logrus"
)

type TokenHandler struct {
	log  *logrus.Logger
	conf dg.AppService
}

func NewTokenHandler(conf dg.AppService, logger *logrus.Logger) TokenHandler {
	utils.InitJWTConfig(conf.Authorization.JWT)
	return TokenHandler{
		log:  logger,
		conf: conf,
	}
}

func (th TokenHandler) JWTValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		respData := utils.ResponseData{
			Status: cg.Fail,
		}

		//List of URL that bypass this JWTValidator middleware
		if req.URL.Path == "/api/v1/renew-token" {
			next.ServeHTTP(res, req)
			return
		}

		authorizationHeader := req.Header.Get("Authorization")
		if !strings.Contains(authorizationHeader, "Bearer") {
			th.log.Error(fmt.Errorf("Invalid Token Format"))
			respData.Message = "Invalid Token Format"
			utils.WriteResponse(res, respData, http.StatusBadRequest)
			return
		}
		accessToken := strings.Replace(authorizationHeader, "Bearer ", "", -1)

		claims, err := utils.CheckAccessToken(accessToken)
		if err != nil {
			respData.Message = "Token expired"
			utils.WriteResponse(res, respData, http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(req.Context(), "session", claims["session"])
		req = req.WithContext(ctx)

		next.ServeHTTP(res, req)
	})
}
