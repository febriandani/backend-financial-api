package user

import (
	"github.com/febriandani/backend-financial-api/domain/model/general"
	"github.com/febriandani/backend-financial-api/service"
	"github.com/sirupsen/logrus"
)

type HandlerUser struct {
	User UserHandler
}

func NewHandlerUser(sv service.Service, conf general.AppService, logger *logrus.Logger) HandlerUser {
	return HandlerUser{
		User: newUserHandler(sv.User, conf, logger),
	}
}
