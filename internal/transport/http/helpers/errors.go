package helpers

import "github.com/internships-backend/test-backend-M0s1ck/internal/transport/http/oapi"

func NewErrorResponse(code oapi.ErrorResponseErrorCode, message string) oapi.ErrorResponse {
	var resp oapi.ErrorResponse
	resp.Error.Code = code
	resp.Error.Message = message
	return resp
}

func NewInternalErrorResponse(message string) oapi.InternalErrorResponse {
	var resp oapi.InternalErrorResponse
	resp.Error.Code = string(oapi.INTERNALERROR)
	resp.Error.Message = message
	return resp
}
