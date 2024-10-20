package category

import (
	"github.com/febriandani/backend-financial-api/infra"
	"github.com/sirupsen/logrus"
)

type DatabaseCategory struct {
	Category Category
}

func NewDatabaseCategory(db *infra.DatabaseList, logger *logrus.Logger) DatabaseCategory {
	return DatabaseCategory{
		Category: newDatabaseCategory(db, logger),
	}
}
