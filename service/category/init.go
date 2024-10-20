package category

import (
	"github.com/febriandani/backend-financial-api/domain/model/general"
	"github.com/febriandani/backend-financial-api/infra"
	rc "github.com/febriandani/backend-financial-api/repository/category"
	"github.com/sirupsen/logrus"
)

type ServiceCategory struct {
	Category Category
}

func NewServiceCategory(database rc.DatabaseCategory, logger *logrus.Logger, conf general.AppService, dbList *infra.DatabaseList) ServiceCategory {
	return ServiceCategory{
		Category: newCategoryService(database, logger, dbList, conf),
	}
}
