package routes

import (
	"net/http"

	"github.com/febriandani/backend-financial-api/domain/model/general"
	"github.com/febriandani/backend-financial-api/handler/api"
	"github.com/gorilla/mux"
)

func getV1(freeRoute, routerJWT *mux.Router, _ *general.AppService, handler api.Handler) {

	freeRoute.HandleFunc("/v1/registration", handler.User.User.RegistrationUser).Methods(http.MethodPost)
	freeRoute.HandleFunc("/v1/login", handler.User.User.LoginUser).Methods(http.MethodPost)
	freeRoute.HandleFunc("/v1/change-password", handler.User.User.UpdatePassword).Methods(http.MethodPost)
	routerJWT.HandleFunc("/v1/transaction", handler.Transaction.Transaction.CreateTransaction).Methods(http.MethodPost)
	routerJWT.HandleFunc("/v1/summary-transaction", handler.Transaction.Transaction.GetSummaryHome).Methods(http.MethodGet)
	routerJWT.HandleFunc("/v1/transactions", handler.Transaction.Transaction.GetTransactions).Methods(http.MethodGet)
	routerJWT.HandleFunc("/v1/category", handler.Category.Category.CreateCategory).Methods(http.MethodPost)
	routerJWT.HandleFunc("/v1/categories", handler.Category.Category.GetCategory).Methods(http.MethodGet)
	routerJWT.HandleFunc("/v1/update-category", handler.Category.Category.UpdateCategory).Methods(http.MethodPost)
	routerJWT.HandleFunc("/v1/delete-category", handler.Category.Category.DeleteCategory).Methods(http.MethodPost)
	routerJWT.HandleFunc("/v1/update-transaction", handler.Transaction.Transaction.UpdateTransaction).Methods(http.MethodPost)
	routerJWT.HandleFunc("/v1/delete-transaction", handler.Transaction.Transaction.DeleteTransaction).Methods(http.MethodPost)
}
