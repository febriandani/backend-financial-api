package transaction

import (
	"github.com/febriandani/backend-financial-api/infra"
	"github.com/sirupsen/logrus"
)

type DatabaseTransaction struct {
	Transaction Transaction
}

func NewDatabaseTransaction(db *infra.DatabaseList, logger *logrus.Logger) DatabaseTransaction {
	return DatabaseTransaction{
		Transaction: newDatabaseTransaction(db, logger),
	}
}
