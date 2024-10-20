package service

import (
	"github.com/febriandani/backend-financial-api/domain/model/general"
	"github.com/febriandani/backend-financial-api/infra"
	"github.com/febriandani/backend-financial-api/repository"
	sc "github.com/febriandani/backend-financial-api/service/category"
	st "github.com/febriandani/backend-financial-api/service/transaction"
	su "github.com/febriandani/backend-financial-api/service/user"
	"github.com/sirupsen/logrus"
)

type Service struct {
	User        su.ServiceUser
	Transaction st.ServiceTransaction
	Category    sc.ServiceCategory
}

func NewService(repo repository.Repo, conf general.AppService, dbList *infra.DatabaseList, logger *logrus.Logger) Service {
	return Service{
		User:        su.NewServiceUser(repo.DatabaseUser, logger, conf, dbList),
		Transaction: st.NewServiceTransaction(repo.DatabaseTransaction, repo.DatabaseCategory, logger, conf, dbList),
		Category:    sc.NewServiceCategory(repo.DatabaseCategory, logger, conf, dbList),
	}
}
