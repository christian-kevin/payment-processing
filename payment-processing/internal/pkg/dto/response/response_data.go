package response

import (
	"encoding/json"
	"net/http"
	"spenmo/payment-processing/payment-processing/internal/pkg/errutil"
)

/*
	For OK Response, err_code is 0
	Add lists of err_code here:
	- 401 Unauthorized
*/

type CommonResponse struct {
	Code       int         `json:"code"`
	ErrMessage string      `json:"err_msg"`
	Data       interface{} `json:"data"`
}

func errorToHTTPStatus(err error) int {
	switch err {
	case errutil.ErrWalletNotFound:
		return http.StatusNotFound
	case errutil.ErrInvalidParam,
		errutil.ErrWalletAlreadyExist,
		errutil.ErrDuplicateRequest:
		return http.StatusBadRequest
	case errutil.ErrUnauthorized,
		errutil.ErrContextValueNotFound:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

func writeFailedEncode(w http.ResponseWriter) {
	errResp := CommonResponse{
		Code:       500,
		ErrMessage: "server error",
		Data:       nil,
	}

	_ = json.NewEncoder(w).Encode(&errResp)
}

func writeResponse(w http.ResponseWriter, resp CommonResponse) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(resp.Code)
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		writeFailedEncode(w)
	}
}

func WriteResponse(w http.ResponseWriter, data interface{}, err error) {
	errString := ""
	if err != nil {
		errString = err.Error()
	}
	resp := CommonResponse{
		Code:       errorToHTTPStatus(err),
		ErrMessage: errString,
	}
	if err == nil {
		resp.Code = 200
		resp.Data = data
	}

	writeResponse(w, resp)
}
