package utils

import (
	"encoding/json"
	"net/http"
)

type ResponseHTTP struct {
	StatusCode int
	Response   ResponseData
}

type ResponseData struct {
	Status  string      `json:"status"`
	Source  string      `json:"source,omitempty"`
	Message string      `json:"message,omitempty"`
	Detail  interface{} `json:"detail"`
}

type ResponseDataV2 struct {
	Status  string            `json:"status"`
	Message map[string]string `json:"message,omitempty"`
	Detail  interface{}       `json:"detail,omitempty"`
}

type Response interface{}

func WriteResponse(res http.ResponseWriter, resp Response, code int) {
	res.Header().Set("Content-Type", "application/json")
	r, _ := json.Marshal(resp)

	res.WriteHeader(code)
	res.Write(r)
	return
}
