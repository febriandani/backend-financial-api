package transaction

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	cg "github.com/febriandani/backend-financial-api/domain/constants/general"
	"github.com/febriandani/backend-financial-api/domain/model/general"
	mt "github.com/febriandani/backend-financial-api/domain/model/transaction"
	"github.com/febriandani/backend-financial-api/domain/utils"
	st "github.com/febriandani/backend-financial-api/service/transaction"
	"github.com/sirupsen/logrus"
)

type TransactionHandler struct {
	transaction st.ServiceTransaction
	conf        general.AppService
	log         *logrus.Logger
}

func newTransactionHandler(transaction st.ServiceTransaction, conf general.AppService, logger *logrus.Logger) TransactionHandler {
	return TransactionHandler{
		transaction: transaction,
		conf:        conf,
		log:         logger,
	}
}

func (th TransactionHandler) CreateTransaction(res http.ResponseWriter, req *http.Request) {
	respData := &utils.ResponseDataV2{
		Status: cg.Fail,
	}
	var param mt.TransactionRequest

	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		respData.Message = map[string]string{
			"en": cg.HandlerErrorRequestDataEmpty,
			"id": cg.HandlerErrorRequestDataEmptyID,
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(reqBody, &param)
	if err != nil {
		respData.Message = map[string]string{
			"en": cg.HandlerErrorRequestDataNotValid,
			"id": cg.HandlerErrorRequestDataNotValidID,
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	dataSession, err := utils.GetUserIDFromToken(fmt.Sprintf("%v", req.Context().Value(cg.SessionContextKey)), th.conf.KeyData.User)
	if err != nil {
		respData.Message = map[string]string{
			"en": th.conf.KeyData.User,
			"id": cg.HandlerErrorAuthInvalidID,
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	var session cg.CredentialData

	err = json.Unmarshal([]byte(dataSession), &session)
	if err != nil {
		respData.Message = map[string]string{
			"en": cg.HandlerErrorAuthInvalid,
			"id": cg.HandlerErrorAuthInvalidID,
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	param.CreatedBy = session.Fullname
	param.UserID = session.ID

	message, err := th.transaction.Transaction.CreateTransaction(req.Context(), param)
	if err != nil {
		respData.Message = message
		utils.WriteResponse(res, respData, http.StatusInternalServerError)
		return
	}

	respData = &utils.ResponseDataV2{
		Status:  cg.Success,
		Message: message,
	}

	utils.WriteResponse(res, respData, http.StatusOK)
	return
}

func (th TransactionHandler) GetSummaryHome(res http.ResponseWriter, req *http.Request) {
	respData := &utils.ResponseDataV2{
		Status: cg.Fail,
	}

	dataSession, err := utils.GetUserIDFromToken(fmt.Sprintf("%v", req.Context().Value(cg.SessionContextKey)), th.conf.KeyData.User)
	if err != nil {
		respData.Message = map[string]string{
			"en": cg.HandlerErrorAuthInvalid,
			"id": cg.HandlerErrorAuthInvalidID,
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	var session cg.CredentialData

	err = json.Unmarshal([]byte(dataSession), &session)
	if err != nil {
		respData.Message = map[string]string{
			"en": cg.HandlerErrorAuthInvalid,
			"id": cg.HandlerErrorAuthInvalidID,
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	data, message, err := th.transaction.Transaction.GetHomeSummaryByUserId(req.Context(), int(session.ID))
	if err != nil {
		respData.Message = message
		utils.WriteResponse(res, respData, http.StatusInternalServerError)
		return
	}

	respData = &utils.ResponseDataV2{
		Status:  cg.Success,
		Message: message,
		Detail:  data,
	}

	utils.WriteResponse(res, respData, http.StatusOK)
	return
}

func (th TransactionHandler) GetTransactions(res http.ResponseWriter, req *http.Request) {
	respData := &utils.ResponseDataV2{
		Status: cg.Fail,
	}

	var param mt.Filter

	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		respData.Message = map[string]string{
			"en": cg.HandlerErrorRequestDataEmpty,
			"id": cg.HandlerErrorRequestDataEmptyID,
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(reqBody, &param)
	if err != nil {
		respData.Message = map[string]string{
			"en": cg.HandlerErrorRequestDataNotValid,
			"id": cg.HandlerErrorRequestDataNotValidID,
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	dataSession, err := utils.GetUserIDFromToken(fmt.Sprintf("%v", req.Context().Value(cg.SessionContextKey)), th.conf.KeyData.User)
	if err != nil {
		respData.Message = map[string]string{
			"en": cg.HandlerErrorAuthInvalid,
			"id": cg.HandlerErrorAuthInvalidID,
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	var session cg.CredentialData

	err = json.Unmarshal([]byte(dataSession), &session)
	if err != nil {
		respData.Message = map[string]string{
			"en": cg.HandlerErrorAuthInvalid,
			"id": cg.HandlerErrorAuthInvalidID,
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	data, message, err := th.transaction.Transaction.GetTransactions(req.Context(), param)
	if err != nil {
		respData.Message = message
		utils.WriteResponse(res, respData, http.StatusInternalServerError)
		return
	}

	respData = &utils.ResponseDataV2{
		Status:  cg.Success,
		Message: message,
		Detail:  data,
	}

	utils.WriteResponse(res, respData, http.StatusOK)
	return
}

func (th TransactionHandler) UpdateTransaction(res http.ResponseWriter, req *http.Request) {
	respData := &utils.ResponseDataV2{
		Status: cg.Fail,
	}

	var param mt.TransactionRequestUpdate

	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		respData.Message = map[string]string{
			"en": cg.HandlerErrorRequestDataEmpty,
			"id": cg.HandlerErrorRequestDataEmptyID,
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(reqBody, &param)
	if err != nil {
		respData.Message = map[string]string{
			"en": cg.HandlerErrorRequestDataNotValid,
			"id": cg.HandlerErrorRequestDataNotValidID,
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	dataSession, err := utils.GetUserIDFromToken(fmt.Sprintf("%v", req.Context().Value(cg.SessionContextKey)), th.conf.KeyData.User)
	if err != nil {
		respData.Message = map[string]string{
			"en": cg.HandlerErrorAuthInvalid,
			"id": cg.HandlerErrorAuthInvalidID,
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	var session cg.CredentialData

	err = json.Unmarshal([]byte(dataSession), &session)
	if err != nil {
		respData.Message = map[string]string{
			"en": cg.HandlerErrorAuthInvalid,
			"id": cg.HandlerErrorAuthInvalidID,
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	param.UpdatedBy = session.Fullname
	param.UserID = int(session.ID)

	message, err := th.transaction.Transaction.UpdateTransaction(req.Context(), param)
	if err != nil {
		respData.Message = message
		utils.WriteResponse(res, respData, http.StatusInternalServerError)
		return
	}

	respData = &utils.ResponseDataV2{
		Status:  cg.Success,
		Message: message,
	}

	utils.WriteResponse(res, respData, http.StatusOK)
	return
}

func (th TransactionHandler) DeleteTransaction(res http.ResponseWriter, req *http.Request) {
	respData := &utils.ResponseDataV2{
		Status: cg.Fail,
	}

	var param mt.TransactionRequestUpdate

	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		respData.Message = map[string]string{
			"en": cg.HandlerErrorRequestDataEmpty,
			"id": cg.HandlerErrorRequestDataEmptyID,
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(reqBody, &param)
	if err != nil {
		respData.Message = map[string]string{
			"en": cg.HandlerErrorRequestDataNotValid,
			"id": cg.HandlerErrorRequestDataNotValidID,
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	dataSession, err := utils.GetUserIDFromToken(fmt.Sprintf("%v", req.Context().Value(cg.SessionContextKey)), th.conf.KeyData.User)
	if err != nil {
		respData.Message = map[string]string{
			"en": cg.HandlerErrorAuthInvalid,
			"id": cg.HandlerErrorAuthInvalidID,
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	var session cg.CredentialData

	err = json.Unmarshal([]byte(dataSession), &session)
	if err != nil {
		respData.Message = map[string]string{
			"en": cg.HandlerErrorAuthInvalid,
			"id": cg.HandlerErrorAuthInvalidID,
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	param.UserID = int(session.ID)

	message, err := th.transaction.Transaction.DeleteTransaction(req.Context(), param)
	if err != nil {
		respData.Message = message
		utils.WriteResponse(res, respData, http.StatusInternalServerError)
		return
	}

	respData = &utils.ResponseDataV2{
		Status:  cg.Success,
		Message: message,
	}

	utils.WriteResponse(res, respData, http.StatusOK)
	return
}
