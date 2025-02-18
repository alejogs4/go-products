package api

import (
	"encoding/json"
	"net/http"
)

type SuccessResponse struct {
	Content interface{} `json:"content"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	AppCode string `json:"app_code"`
}

const (
	InvalidRequestCode      = "INVALID_REQUEST"
	InternalServerErrorCode = "INTERNAL_SERVER_ERROR"
	MethodNotAllowedCode    = "METHOD_NOT_ALLOWED"
)

func Success(response http.ResponseWriter, data interface{}) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(response).Encode(SuccessResponse{Content: data})
}

func InternalServerError(response http.ResponseWriter, message string) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(response).Encode(ErrorResponse{Message: message, AppCode: InternalServerErrorCode})
}

func InvalidRequest(response http.ResponseWriter, message string) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusBadRequest)
	_ = json.NewEncoder(response).Encode(ErrorResponse{Message: message, AppCode: InvalidRequestCode})
}

func methodNotAllowed(response http.ResponseWriter) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusMethodNotAllowed)
	_ = json.NewEncoder(response).Encode(ErrorResponse{Message: "Method not allowed", AppCode: MethodNotAllowedCode})
}
