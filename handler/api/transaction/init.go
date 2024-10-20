package transaction

import (
	"github.com/febriandani/backend-financial-api/domain/model/general"
	"github.com/febriandani/backend-financial-api/service"
	"github.com/sirupsen/logrus"
)

type HandlerTransaction struct {
	Transaction TransactionHandler
}

func NewHandlerTransaction(sv service.Service, conf general.AppService, logger *logrus.Logger) HandlerTransaction {
	return HandlerTransaction{
		Transaction: newTransactionHandler(sv.Transaction, conf, logger),
	}
}
