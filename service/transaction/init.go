package transaction

import (
	"github.com/febriandani/backend-financial-api/domain/model/general"
	"github.com/febriandani/backend-financial-api/infra"
	rc "github.com/febriandani/backend-financial-api/repository/category"
	rt "github.com/febriandani/backend-financial-api/repository/transaction"
	"github.com/sirupsen/logrus"
)

type ServiceTransaction struct {
	Transaction Transaction
}

func NewServiceTransaction(database rt.DatabaseTransaction, databaseCategory rc.DatabaseCategory, logger *logrus.Logger, conf general.AppService, dbList *infra.DatabaseList) ServiceTransaction {
	return ServiceTransaction{
		Transaction: newTransactionService(database, databaseCategory, logger, dbList, conf),
	}
}
