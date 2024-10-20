package category

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	cg "github.com/febriandani/backend-financial-api/domain/constants/general"
	mt "github.com/febriandani/backend-financial-api/domain/model/category"
	"github.com/febriandani/backend-financial-api/domain/model/general"
	"github.com/febriandani/backend-financial-api/domain/utils"
	st "github.com/febriandani/backend-financial-api/service/category"
	"github.com/sirupsen/logrus"
)

type CategoryHandler struct {
	category st.ServiceCategory
	conf     general.AppService
	log      *logrus.Logger
}

func newCategoryHandler(category st.ServiceCategory, conf general.AppService, logger *logrus.Logger) CategoryHandler {
	return CategoryHandler{
		category: category,
		conf:     conf,
		log:      logger,
	}
}

func (th CategoryHandler) CreateCategory(res http.ResponseWriter, req *http.Request) {
	respData := &utils.ResponseDataV2{
		Status: cg.Fail,
	}
	var param mt.CategoryRequest

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
	param.UserID = int(session.ID)

	message, err := th.category.Category.CreateCategory(req.Context(), param)
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

func (th CategoryHandler) GetCategory(res http.ResponseWriter, req *http.Request) {
	respData := &utils.ResponseDataV2{
		Status: cg.Fail,
	}

	categoryType := req.URL.Query().Get("type")

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

	data, message, err := th.category.Category.GetCategoryByUserId(req.Context(), int(session.ID), categoryType)
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

func (th CategoryHandler) UpdateCategory(res http.ResponseWriter, req *http.Request) {
	respData := &utils.ResponseDataV2{
		Status: cg.Fail,
	}
	var param mt.CategoryRequestUpdate

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

	param.UpdatedBy = session.Fullname
	param.UserID = int(session.ID)

	message, err := th.category.Category.UpdateCategory(req.Context(), param)
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

func (th CategoryHandler) DeleteCategory(res http.ResponseWriter, req *http.Request) {
	respData := &utils.ResponseDataV2{
		Status: cg.Fail,
	}
	var param mt.CategoryRequestUpdate

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

	param.UserID = int(session.ID)

	message, err := th.category.Category.DeleteCategory(req.Context(), param)
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
