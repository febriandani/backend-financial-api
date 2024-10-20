package repository

import (
	"github.com/febriandani/backend-financial-api/infra"
	rc "github.com/febriandani/backend-financial-api/repository/category"
	rt "github.com/febriandani/backend-financial-api/repository/transaction"
	ru "github.com/febriandani/backend-financial-api/repository/user"
	"github.com/sirupsen/logrus"
)

type Repo struct {
	DatabaseUser        ru.DatabaseUser
	DatabaseTransaction rt.DatabaseTransaction
	DatabaseCategory    rc.DatabaseCategory
}

func NewRepo(database *infra.DatabaseList, logger *logrus.Logger) Repo {
	return Repo{
		DatabaseUser:        ru.NewDatabaseUser(database, logger),
		DatabaseTransaction: rt.NewDatabaseTransaction(database, logger),
		DatabaseCategory:    rc.NewDatabaseCategory(database, logger),
	}
}
