package api

import (
	"github.com/febriandani/backend-financial-api/domain/model/general"
	"github.com/febriandani/backend-financial-api/handler/api/authorization"
	hc "github.com/febriandani/backend-financial-api/handler/api/category"
	ht "github.com/febriandani/backend-financial-api/handler/api/transaction"
	hu "github.com/febriandani/backend-financial-api/handler/api/user"
	"github.com/febriandani/backend-financial-api/service"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Token       authorization.TokenHandler
	User        hu.HandlerUser
	Transaction ht.HandlerTransaction
	Category    hc.HandlerCategory
}

func NewHandler(sv service.Service, conf general.AppService, logger *logrus.Logger) Handler {
	return Handler{
		Token:       authorization.NewTokenHandler(conf, logger),
		User:        hu.NewHandlerUser(sv, conf, logger),
		Transaction: ht.NewHandlerTransaction(sv, conf, logger),
		Category:    hc.NewHandlerCategory(sv, conf, logger),
	}
}
