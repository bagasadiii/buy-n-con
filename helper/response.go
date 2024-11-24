package helper

import (
	"encoding/json"
	"net/http"
)


type Response struct {
	Status  int				`json:"status"`
	Message string			`json:"message"`
	Data    interface{}		`json:"data"`
	Err		interface{}		`json:"err"`
}

func JSONResponse(w http.ResponseWriter, status int, data interface{}){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
func UnauthorizedErr(msg string, err error)*Response{
	ErrMsg(err, msg)
	return &Response{
		Status: http.StatusUnauthorized,
		Message: msg,
		Data: nil,
		Err: err,
	}
}
func BadRequestErr(msg string, err error)*Response{
	ErrMsg(err, msg)
	return &Response{
		Status: http.StatusBadRequest,
		Message: msg,
		Data: nil,
		Err: err,
	}
}
func InternalErr(msg string, err error)*Response{
	ErrMsg(err, msg)
	return &Response{
		Status: http.StatusInternalServerError,
		Message: msg,
		Data: nil,
		Err: err,
	}
}
func ForbiddenErr(msg string, err error)*Response{
	ErrMsg(err, msg)
	return &Response{
		Status: http.StatusForbidden,
		Message: msg,
		Data: nil,
		Err: err,
	}
}