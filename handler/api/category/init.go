package category

import (
	"github.com/febriandani/backend-financial-api/domain/model/general"
	"github.com/febriandani/backend-financial-api/service"
	"github.com/sirupsen/logrus"
)

type HandlerCategory struct {
	Category CategoryHandler
}

func NewHandlerCategory(sv service.Service, conf general.AppService, logger *logrus.Logger) HandlerCategory {
	return HandlerCategory{
		Category: newCategoryHandler(sv.Category, conf, logger),
	}
}
