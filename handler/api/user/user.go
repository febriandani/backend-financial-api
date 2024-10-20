package user

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	cg "github.com/febriandani/backend-financial-api/domain/constants/general"
	"github.com/febriandani/backend-financial-api/domain/model/general"
	mu "github.com/febriandani/backend-financial-api/domain/model/user"
	"github.com/febriandani/backend-financial-api/domain/utils"
	su "github.com/febriandani/backend-financial-api/service/user"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	user su.ServiceUser
	conf general.AppService
	log  *logrus.Logger
}

func newUserHandler(user su.ServiceUser, conf general.AppService, logger *logrus.Logger) UserHandler {
	return UserHandler{
		user: user,
		conf: conf,
		log:  logger,
	}
}

func (uh UserHandler) RegistrationUser(res http.ResponseWriter, req *http.Request) {
	respData := &utils.ResponseDataV2{
		Status: cg.Fail,
	}
	var param mu.RegistrationUser

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

	message, err := uh.user.User.CreateRegistrationUser(req.Context(), param)
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

func (uh UserHandler) UpdatePassword(res http.ResponseWriter, req *http.Request) {
	respData := &utils.ResponseDataV2{
		Status: cg.Fail,
	}
	var param mu.ForgotPasswordRequest

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

	message, err := uh.user.User.ForgotPassword(req.Context(), param)
	if err != nil {
		if err.Error() == "PS-011" {
			respData.Message = message
			respData.Status = "PS-011"
			utils.WriteResponse(res, respData, http.StatusNotAcceptable)
			return
		} else {
			respData.Message = message
			utils.WriteResponse(res, respData, http.StatusInternalServerError)
			return
		}

	}

	respData = &utils.ResponseDataV2{
		Status:  cg.Success,
		Message: message,
	}

	utils.WriteResponse(res, respData, http.StatusOK)
	return
}

func (uh UserHandler) LoginUser(res http.ResponseWriter, req *http.Request) {
	respData := &utils.ResponseDataV2{
		Status: cg.Fail,
	}
	var param mu.LoginRequest

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

	data, message, err := uh.user.User.Login(req.Context(), param)
	if err != nil {
		if err.Error() == "FailedServer" {
			respData.Message = message
			utils.WriteResponse(res, respData, http.StatusInternalServerError)
			return
		} else if err.Error() == "FailedPassword" {
			respData.Message = message
			respData.Status = "FailedPassword"
			utils.WriteResponse(res, respData, http.StatusForbidden)
			return
		} else if err.Error() == "UserNA" {
			respData.Message = message
			respData.Status = "FailedUserNA"
			utils.WriteResponse(res, respData, http.StatusNotFound)
			return
		} else {
			respData.Message = message
			utils.WriteResponse(res, respData, http.StatusInternalServerError)
			return
		}

	}

	respData = &utils.ResponseDataV2{
		Status:  cg.Success,
		Message: message,
		Detail:  data,
	}

	utils.WriteResponse(res, respData, http.StatusOK)
	return
}
