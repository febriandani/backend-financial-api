package user

import (
	"github.com/febriandani/backend-financial-api/domain/model/general"
	"github.com/febriandani/backend-financial-api/infra"
	ru "github.com/febriandani/backend-financial-api/repository/user"
	"github.com/sirupsen/logrus"
)

type ServiceUser struct {
	User User
}

func NewServiceUser(database ru.DatabaseUser, logger *logrus.Logger, conf general.AppService, dbList *infra.DatabaseList) ServiceUser {
	return ServiceUser{
		User: newUserService(database, logger, dbList, conf),
	}
}
